export interface User {
  id: number
  firstName: string
  lastName: string
  email: string
  role: "user" | "manager" | "admin"
  image?: string
  blocked: boolean
  createdAt: string
  updatedAt: string
}

export interface Pet {
  id: number
  name: string
  description?: string
  price: number
  breed: string
  age: number
  gender: "male" | "female"
  sterilized: boolean
  image?: string
  ownerId: number
  createdAt: string
  updatedAt: string
}

export interface Product {
  id: number
  name: string
  description?: string
  price: number
  stock: number
  category: string
  brand?: string
  image?: string
  mass: number
  ownerId: number
  createdAt: string
  updatedAt: string
}

export interface Stats {
  users: number
  totalPets: number
  ownedPets: number
  storePets: number
  totalProducts: number
  ownedProducts: number
  storeProducts: number
}

export interface AuthResponse {
  token: string
  user: User
}

export interface ApiError {
  error: string
}
