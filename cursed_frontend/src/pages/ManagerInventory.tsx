// src/pages/ManagerInventory.tsx (modified: unified nav, no sidebar, full-width content, adjusted loading)
"use client"

import { useState, useEffect } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/Button"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/Input"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Checkbox } from "@/components/ui/checkbox"
import { Card, CardContent } from "@/components/ui/Card"
import { Plus, Edit, Trash2 } from "lucide-react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import type { Pet, Product } from "@/types/types.ts"
import type { ApiResponse } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"
import { AppNav } from "@/components/AppNav.tsx";  // Unified nav

const petSchema = z.object({
    name: z.string().min(1, "Name is required"),
    description: z.string().optional(),
    price: z.coerce.number().positive("Price must be positive"),
    breed: z.string().min(1, "Breed is required"),
    age: z.coerce.number().min(0, "Age must be non-negative"),
    gender: z.enum(["male", "female"], { required_error: "Gender is required" }),
    sterilized: z.boolean().default(false),
    image: z.string().url("Invalid URL").optional(),
    ownerId: z.coerce.number().min(0).optional().default(0),
})

const productSchema = z.object({
    name: z.string().min(1, "Name is required"),
    description: z.string().optional(),
    price: z.coerce.number().positive("Price must be positive"),
    stock: z.coerce.number().min(0, "Stock must be non-negative"),
    category: z.string().min(1, "Category is required"),
    brand: z.string().optional(),
    image: z.string().url("Invalid URL").optional(),
    mass: z.coerce.number().min(0, "Mass must be non-negative"),
    ownerId: z.coerce.number().min(0).optional().default(0),
})

type PetFormData = z.infer<typeof petSchema>
type ProductFormData = z.infer<typeof productSchema>

export function ManagerInventory() {
    const [activeTab, setActiveTab] = useState<"pets" | "products">("pets")
    const [pets, setPets] = useState<Pet[]>([])
    const [products, setProducts] = useState<Product[]>([])
    const [loading, setLoading] = useState(true)
    const [editingId, setEditingId] = useState<number | null>(null)
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [deletingId, setDeletingId] = useState<number | null>(null)

    useEffect(() => {
        fetchItems()
    }, [activeTab])

    const fetchItems = async () => {
        try {
            setLoading(true)
            const [petsRes, productsRes] = await Promise.all([
                api.get<ApiResponse<Pet[]>>("/pets"),
                api.get<ApiResponse<Product[]>>("/products"),
            ]);
            setPets(petsRes.data.data || []);  // {success, data: []}
            setProducts(productsRes.data.data || []);
        } catch (error) {
            toast.error("Failed to load inventory")
        } finally {
            setLoading(false)
        }
    }

    // Pet Form
    const petForm = useForm<PetFormData>({
        resolver: zodResolver(petSchema),
        defaultValues: {
            name: "",
            description: "",
            price: 0,
            breed: "",
            age: 0,
            gender: "male",
            sterilized: false,
            image: "",
            ownerId: 0,
        },
    })

    // Product Form — FIXED: useForm<ProductFormData> instead of <Product>
    const productForm = useForm<ProductFormData>({
        resolver: zodResolver(productSchema),
        defaultValues: {
            name: "",
            description: "",
            price: 0,
            stock: 0,
            category: "",
            brand: "",
            image: "",
            mass: 0,
            ownerId: 0,
        },
    })

    const resetForms = () => {
        petForm.reset()
        productForm.reset()
        setEditingId(null)
        setIsModalOpen(false)
    }

    const openEditModal = (item: Pet | Product) => {
        setEditingId(item.id)
        if ("breed" in item) {
            // Pet
            petForm.reset({
                name: item.name,
                description: item.description || "",
                price: item.price,
                breed: item.breed,
                age: item.age,
                gender: item.gender,
                sterilized: item.sterilized,
                image: item.image || "",
                ownerId: item.ownerId,
            })
        } else {
            // Product
            productForm.reset({
                name: item.name,
                description: item.description || "",
                price: item.price,
                stock: item.stock,
                category: item.category,
                brand: item.brand || "",
                image: item.image || "",
                mass: item.mass,
                ownerId: item.ownerId,
            })
        }
        setIsModalOpen(true)
    }

    const handleSubmitPet = async (data: PetFormData) => {
        try {
            if (editingId) {
                // Update
                await api.put<ApiResponse<Pet>>(`/pets/${editingId}`, data)
                toast.success("Pet updated successfully")
            } else {
                // Create
                await api.post<ApiResponse<Pet>>("/pets", data)
                toast.success("Pet created successfully")
            }
            fetchItems()
            resetForms()
        } catch (error: any) {
            toast.error(error.response?.data?.error || "Failed to save pet")
        }
    }

    const handleSubmitProduct = async (data: ProductFormData) => {
        try {
            if (editingId) {
                // Update
                await api.put<ApiResponse<Product>>(`/products/${editingId}`, data)
                toast.success("Product updated successfully")
            } else {
                // Create
                await api.post<ApiResponse<Product>>("/products", data)
                toast.success("Product created successfully")
            }
            fetchItems()
            resetForms()
        } catch (error: any) {
            toast.error(error.response?.data?.error || "Failed to save product")
        }
    }

    const confirmDelete = (id: number, type: "pet" | "product") => {
        setDeletingId(id)
        toast(
            (t) => (
                <div className="flex justify-between items-center">
                    <span>Delete this {type}?</span>
                    <div className="flex gap-2">
                        <Button
                            size="sm"
                            variant="outline"
                            onClick={() => {
                                toast.dismiss(t.id)
                                setDeletingId(null)
                            }}
                        >
                            Cancel
                        </Button>
                        <Button
                            size="sm"
                            variant="destructive"
                            onClick={async () => {
                                await handleDelete(id, type)
                                toast.dismiss(t.id)
                                setDeletingId(null)
                            }}
                        >
                            Delete
                        </Button>
                    </div>
                </div>
            ),
            { duration: Infinity, position: "top-center" }
        )
    }

    const handleDelete = async (id: number, type: "pet" | "product") => {
        try {
            await api.delete<ApiResponse<null>>(`/${type}s/${id}`)
            toast.success(`${type.charAt(0).toUpperCase() + type.slice(1)} deleted successfully`)
            fetchItems()
        } catch (error: any) {
            toast.error("Delete failed")
        }
    }

    if (loading) {
        return (
            <div className="min-h-screen bg-background">
                <AppNav />
                <div className="flex h-[70vh] items-center justify-center">
                    <div className="animate-spin rounded-full border-2 border-primary h-8 w-8" />
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-background">
            <AppNav />
            <main className="overflow-y-auto">
                <div className="p-8">
                    <div className="mb-6 flex items-center justify-between">
                        <h1 className="text-3xl font-bold">Inventory Management</h1>
                        <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
                            <DialogTrigger asChild>
                                {/* FIXED: Added setIsModalOpen(true) to ensure modal opens for add */}
                                <Button
                                    onClick={() => {
                                        setEditingId(null)
                                        setIsModalOpen(true)  // Explicitly open modal
                                    }}
                                >
                                    <Plus className="mr-2 h-4 w-4" />
                                    Add {activeTab === "pets" ? "Pet" : "Product"}
                                </Button>
                            </DialogTrigger>
                            <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
                                <DialogHeader>
                                    <DialogTitle>
                                        {editingId ? "Edit" : "Add"} {activeTab === "pets" ? "Pet" : "Product"}
                                    </DialogTitle>
                                </DialogHeader>
                                <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as "pets" | "products")}>
                                    <TabsList>
                                        <TabsTrigger value="pets">Pets</TabsTrigger>
                                        <TabsTrigger value="products">Products</TabsTrigger>
                                    </TabsList>
                                    <TabsContent value="pets" className="mt-4">
                                        <Form {...petForm}>
                                            <form onSubmit={petForm.handleSubmit(handleSubmitPet)} className="space-y-4">
                                                <FormField
                                                    control={petForm.control}
                                                    name="name"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Name</FormLabel>
                                                            <FormControl>
                                                                <Input {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={petForm.control}
                                                    name="breed"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Breed</FormLabel>
                                                            <FormControl>
                                                                <Input {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <div className="grid grid-cols-2 gap-4">
                                                    <FormField
                                                        control={petForm.control}
                                                        name="age"
                                                        render={({ field }) => (
                                                            <FormItem>
                                                                <FormLabel>Age</FormLabel>
                                                                <FormControl>
                                                                    <Input type="number" {...field} />
                                                                </FormControl>
                                                                <FormMessage />
                                                            </FormItem>
                                                        )}
                                                    />
                                                    <FormField
                                                        control={petForm.control}
                                                        name="gender"
                                                        render={({ field }) => (
                                                            <FormItem>
                                                                <FormLabel>Gender</FormLabel>
                                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                                    <FormControl>
                                                                        <SelectTrigger>
                                                                            <SelectValue />
                                                                        </SelectTrigger>
                                                                    </FormControl>
                                                                    <SelectContent>
                                                                        <SelectItem value="male">Male</SelectItem>
                                                                        <SelectItem value="female">Female</SelectItem>
                                                                    </SelectContent>
                                                                </Select>
                                                                <FormMessage />
                                                            </FormItem>
                                                        )}
                                                    />
                                                </div>
                                                <FormField
                                                    control={petForm.control}
                                                    name="price"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Price ($)</FormLabel>
                                                            <FormControl>
                                                                <Input type="number" step="0.01" {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={petForm.control}
                                                    name="description"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Description</FormLabel>
                                                            <FormControl>
                                                                <Textarea {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={petForm.control}
                                                    name="image"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Image URL</FormLabel>
                                                            <FormControl>
                                                                <Input placeholder="https://example.com/image.jpg" {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={petForm.control}
                                                    name="sterilized"
                                                    render={({ field }) => (
                                                        <FormItem className="flex items-center space-x-2">
                                                            <FormControl>
                                                                <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                                                            </FormControl>
                                                            <div className="space-y-1 leading-none">
                                                                <FormLabel>Sterilized</FormLabel>
                                                            </div>
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={petForm.control}
                                                    name="ownerId"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Owner ID (0 for store)</FormLabel>
                                                            <FormControl>
                                                                <Input type="number" {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <Button type="submit" className="w-full">
                                                    {editingId ? "Update" : "Create"} Pet
                                                </Button>
                                            </form>
                                        </Form>
                                    </TabsContent>
                                    <TabsContent value="products" className="mt-4">
                                        <Form {...productForm}>
                                            <form onSubmit={productForm.handleSubmit(handleSubmitProduct)} className="space-y-4">
                                                <FormField
                                                    control={productForm.control}
                                                    name="name"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Name</FormLabel>
                                                            <FormControl>
                                                                <Input {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <div className="grid grid-cols-2 gap-4">
                                                    <FormField
                                                        control={productForm.control}
                                                        name="category"
                                                        render={({ field }) => (
                                                            <FormItem>
                                                                <FormLabel>Category</FormLabel>
                                                                <FormControl>
                                                                    <Input {...field} />
                                                                </FormControl>
                                                                <FormMessage />
                                                            </FormItem>
                                                        )}
                                                    />
                                                    <FormField
                                                        control={productForm.control}
                                                        name="brand"
                                                        render={({ field }) => (
                                                            <FormItem>
                                                                <FormLabel>Brand</FormLabel>
                                                                <FormControl>
                                                                    <Input {...field} />
                                                                </FormControl>
                                                                <FormMessage />
                                                            </FormItem>
                                                        )}
                                                    />
                                                </div>
                                                <div className="grid grid-cols-2 gap-4">
                                                    <FormField
                                                        control={productForm.control}
                                                        name="stock"
                                                        render={({ field }) => (
                                                            <FormItem>
                                                                <FormLabel>Stock</FormLabel>
                                                                <FormControl>
                                                                    <Input type="number" {...field} />
                                                                </FormControl>
                                                                <FormMessage />
                                                            </FormItem>
                                                        )}
                                                    />
                                                    <FormField
                                                        control={productForm.control}
                                                        name="mass"
                                                        render={({ field }) => (
                                                            <FormItem>
                                                                <FormLabel>Mass (kg)</FormLabel>
                                                                <FormControl>
                                                                    <Input type="number" step="0.01" {...field} />
                                                                </FormControl>
                                                                <FormMessage />
                                                            </FormItem>
                                                        )}
                                                    />
                                                </div>
                                                <FormField
                                                    control={productForm.control}
                                                    name="price"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Price ($)</FormLabel>
                                                            <FormControl>
                                                                <Input type="number" step="0.01" {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={productForm.control}
                                                    name="description"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Description</FormLabel>
                                                            <FormControl>
                                                                <Textarea {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={productForm.control}
                                                    name="image"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Image URL</FormLabel>
                                                            <FormControl>
                                                                <Input placeholder="https://example.com/image.jpg" {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <FormField
                                                    control={productForm.control}
                                                    name="ownerId"
                                                    render={({ field }) => (
                                                        <FormItem>
                                                            <FormLabel>Owner ID (0 for store)</FormLabel>
                                                            <FormControl>
                                                                <Input type="number" {...field} />
                                                            </FormControl>
                                                            <FormMessage />
                                                        </FormItem>
                                                    )}
                                                />
                                                <Button type="submit" className="w-full">
                                                    {editingId ? "Update" : "Create"} Product
                                                </Button>
                                            </form>
                                        </Form>
                                    </TabsContent>
                                </Tabs>
                            </DialogContent>
                        </Dialog>
                    </div>

                    <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as "pets" | "products")}>
                        <TabsList>
                            <TabsTrigger value="pets">Pets ({pets.length})</TabsTrigger>
                            <TabsTrigger value="products">Products ({products.length})</TabsTrigger>
                        </TabsList>
                        <TabsContent value="pets" className="mt-6">
                            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                                {pets.map((pet) => (
                                    <Card key={pet.id}>
                                        <img
                                            src={pet.image || "/placeholder-pet.jpg"}
                                            alt={pet.name}
                                            className="aspect-[4/3] w-full rounded-t-lg object-cover"
                                        />
                                        <CardContent className="p-4">
                                            <h3 className="font-semibold">{pet.name}</h3>
                                            <p className="text-sm text-muted-foreground">{pet.breed} • {pet.age}y • {pet.gender}</p>
                                            <p className="mt-2 text-sm text-muted-foreground">${pet.price.toFixed(2)}</p>
                                            <div className="mt-4 flex justify-end space-x-2">
                                                <Button variant="outline" size="sm" onClick={() => openEditModal(pet)}>
                                                    <Edit className="mr-1 h-4 w-4" />
                                                    Edit
                                                </Button>
                                                <Button
                                                    variant="destructive"
                                                    size="sm"
                                                    onClick={() => confirmDelete(pet.id, "pet")}
                                                    disabled={deletingId === pet.id}
                                                >
                                                    <Trash2 className="mr-1 h-4 w-4" />
                                                    Delete
                                                </Button>
                                            </div>
                                        </CardContent>
                                    </Card>
                                ))}
                            </div>
                        </TabsContent>
                        <TabsContent value="products" className="mt-6">
                            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                                {products.map((product) => (
                                    <Card key={product.id}>
                                        <img
                                            src={product.image || "/placeholder-product.jpg"}
                                            alt={product.name}
                                            className="aspect-square w-full rounded-t-lg object-cover"
                                        />
                                        <CardContent className="p-4">
                                            <h3 className="font-semibold">{product.name}</h3>
                                            <p className="text-sm text-muted-foreground">{product.category} • Stock: {product.stock}</p>
                                            <p className="mt-2 text-sm text-muted-foreground">${product.price.toFixed(2)}</p>
                                            <div className="mt-4 flex justify-end space-x-2">
                                                <Button variant="outline" size="sm" onClick={() => openEditModal(product)}>
                                                    <Edit className="mr-1 h-4 w-4" />
                                                    Edit
                                                </Button>
                                                <Button
                                                    variant="destructive"
                                                    size="sm"
                                                    onClick={() => confirmDelete(product.id, "product")}
                                                    disabled={deletingId === product.id}
                                                >
                                                    <Trash2 className="mr-1 h-4 w-4" />
                                                    Delete
                                                </Button>
                                            </div>
                                        </CardContent>
                                    </Card>
                                ))}
                            </div>
                        </TabsContent>
                    </Tabs>
                </div>
            </main>
        </div>
    )
}