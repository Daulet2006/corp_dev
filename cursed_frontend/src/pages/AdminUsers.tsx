// src/pages/AdminUsers.tsx (modified: unified nav, no sidebar, full-width content, adjusted loading)
"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/Button"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Input } from "@/components/ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent } from "@/components/ui/Card"
import { Edit, ShieldAlert, ShieldCheck, UserCog } from "lucide-react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import type { User } from "@/types/types.ts"
import type { ApiResponse } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"
import { AppNav } from "@/components/AppNav.tsx"  // Unified nav
import { useAuthStore } from "@/stores/authStore"

const userSchema = z.object({
    firstName: z.string().min(1, "First name is required"),
    lastName: z.string().min(1, "Last name is required"),
    email: z.string().email("Invalid email"),
    image: z.string().url("Invalid URL").optional(),
})

const roleSchema = z.object({
    role: z.enum(["user", "manager", "admin"], { required_error: "Role is required" }),
})

type UserFormData = z.infer<typeof userSchema>
type RoleFormData = z.infer<typeof roleSchema>

export function AdminUsers() {
    const { user: currentUser } = useAuthStore() // FIXED: Get current user to prevent self-block
    const [users, setUsers] = useState<User[]>([])
    const [loading, setLoading] = useState(true)
    const [editingId, setEditingId] = useState<number | null>(null)
    const [changingRoleId, setChangingRoleId] = useState<number | null>(null)
    const [isEditModalOpen, setIsEditModalOpen] = useState(false)
    const [isRoleModalOpen, setIsRoleModalOpen] = useState(false)

    useEffect(() => {
        fetchUsers()
    }, [])

    const fetchUsers = async () => {
        try {
            setLoading(true)
            const response = await api.get<ApiResponse<User[]>>("/admin/users")
            setUsers(response.data.data || [])
        } catch (error) {
            toast.error("Failed to load users")
        } finally {
            setLoading(false)
        }
    }

    // Edit User Form
    const userForm = useForm<UserFormData>({
        resolver: zodResolver(userSchema),
        defaultValues: {
            firstName: "",
            lastName: "",
            email: "",
            image: "",
        },
    })

    // Change Role Form
    const roleForm = useForm<RoleFormData>({
        resolver: zodResolver(roleSchema),
        defaultValues: {
            role: "user",
        },
    })

    const resetForms = () => {
        userForm.reset()
        roleForm.reset()
        setEditingId(null)
        setIsEditModalOpen(false)
        setChangingRoleId(null)
        setIsRoleModalOpen(false)
    }

    const openEditModal = (user: User) => {
        setEditingId(user.id)
        userForm.reset({
            firstName: user.firstName,
            lastName: user.lastName,
            email: user.email,
            image: user.image || "",
        })
        setIsEditModalOpen(true)
    }

    const openChangeRoleModal = (user: User) => {
        setChangingRoleId(user.id)
        roleForm.reset({ role: user.role })
        setIsRoleModalOpen(true)
    }

    const handleSubmitUser = async (data: UserFormData) => {
        try {
            if (editingId) {
                await api.put<ApiResponse<User>>(`/admin/users/${editingId}`, data)
                toast.success("User updated successfully")
            }
            fetchUsers()
            resetForms()
        } catch (error: any) {
            toast.error(error.response?.data?.error || "Failed to save user")
        }
    }

    const handleSubmitRole = async (data: RoleFormData) => {
        try {
            if (changingRoleId) {
                await api.put<ApiResponse<User>>(`/admin/users/${changingRoleId}/role`, data)
                toast.success("Role changed successfully")
            }
            fetchUsers()
            resetForms()
        } catch (error: any) {
            toast.error(error.response?.data?.error || "Failed to change role")
        }
    }

    const handleToggleBlock = async (user: User) => {
        // FIXED: Prevent self-block
        if (user.id === currentUser?.id) {
            toast.error("Cannot block or unblock yourself")
            return
        }

        try {
            const endpoint = user.blocked ? "unblock" : "block"
            await api.post<ApiResponse<User>>(`/admin/users/${user.id}/${endpoint}`)
            toast.success(`User ${user.blocked ? "unblocked" : "blocked"} successfully`)
            fetchUsers()
        } catch (error: any) {
            toast.error("Toggle failed")
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
                    <h1 className="text-3xl font-bold mb-6">User Management</h1>

                    {/* Edit User Modal */}
                    <Dialog open={isEditModalOpen} onOpenChange={setIsEditModalOpen}>
                        <DialogContent className="max-w-md">
                            <DialogHeader>
                                <DialogTitle>Edit User</DialogTitle>
                            </DialogHeader>
                            <Form {...userForm}>
                                <form onSubmit={userForm.handleSubmit(handleSubmitUser)} className="space-y-4">
                                    <FormField
                                        control={userForm.control}
                                        name="firstName"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>First Name</FormLabel>
                                                <FormControl>
                                                    <Input {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={userForm.control}
                                        name="lastName"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Last Name</FormLabel>
                                                <FormControl>
                                                    <Input {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={userForm.control}
                                        name="email"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Email</FormLabel>
                                                <FormControl>
                                                    <Input type="email" {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={userForm.control}
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
                                    <Button type="submit" className="w-full">Update User</Button>
                                </form>
                            </Form>
                        </DialogContent>
                    </Dialog>

                    {/* Change Role Modal */}
                    <Dialog open={isRoleModalOpen} onOpenChange={setIsRoleModalOpen}>
                        <DialogContent className="max-w-md">
                            <DialogHeader>
                                <DialogTitle>Change Role</DialogTitle>
                            </DialogHeader>
                            <Form {...roleForm}>
                                <form onSubmit={roleForm.handleSubmit(handleSubmitRole)} className="space-y-4">
                                    <FormField
                                        control={roleForm.control}
                                        name="role"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Role</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem value="user">User</SelectItem>
                                                        <SelectItem value="manager">Manager</SelectItem>
                                                        <SelectItem value="admin">Admin</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <Button type="submit" className="w-full">Change Role</Button>
                                </form>
                            </Form>
                        </DialogContent>
                    </Dialog>

                    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                        {users.map((user) => (
                            <Card key={user.id}>
                                <CardContent className="p-4">
                                    <div className="flex items-center justify-between mb-2">
                                        <h3 className="font-semibold">{user.firstName} {user.lastName}</h3>
                                        {user.blocked ? (
                                            <ShieldAlert className="h-5 w-5 text-destructive" />
                                        ) : (
                                            <ShieldCheck className="h-5 w-5 text-success" />
                                        )}
                                    </div>
                                    <p className="text-sm text-muted-foreground mb-1">{user.email}</p>
                                    <p className="text-sm font-medium mb-4 capitalize">{user.role}</p>
                                    <div className="flex justify-end space-x-2">
                                        <Button variant="outline" size="sm" onClick={() => openEditModal(user)}>
                                            <Edit className="mr-1 h-4 w-4" />
                                            Edit
                                        </Button>
                                        <Button variant="outline" size="sm" onClick={() => openChangeRoleModal(user)}>
                                            <UserCog className="mr-1 h-4 w-4" />
                                            Role
                                        </Button>
                                        <Button
                                            variant={user.blocked ? "default" : "outline"}
                                            size="sm"
                                            onClick={() => handleToggleBlock(user)}
                                            disabled={user.id === currentUser?.id} // FIXED: Disable for self
                                        >
                                            {user.blocked ? "Unblock" : "Block"}
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </div>
            </main>
        </div>
    )
}