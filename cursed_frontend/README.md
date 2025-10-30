# ZooPet Store - React + Vite + TypeScript

A modern, full-featured zoo pet store application with distinct interfaces for three user roles: Customer, Manager, and Admin.

## Features

### Multi-Role Interfaces
- **Customer Interface**: Browse and adopt pets, shop products with warm and inviting design
- **Manager Interface**: Inventory management with data-focused layout (teal accent color)
- **Admin Interface**: User management and analytics dashboard (amber primary color)

### Core Features
- 🌓 Light/Dark theme switching
- 🌍 Multi-language support (EN/RU/KZ)
- 📱 Fully responsive design
- 🎨 Premium, trustworthy design aesthetic
- 🔐 JWT authentication with role-based access
- 🛒 Pet and product purchasing system
- 📊 Admin analytics dashboard
- 📦 Manager inventory management

### Technology Stack
- React 18
- Vite
- TypeScript
- Tailwind CSS
- React Router v6
- Zustand (state management)
- Axios (API calls)
- React Hook Form + Zod (form validation)
- React Hot Toast (notifications)
- Lucide React (icons)

## Getting Started

1. Install dependencies:
\`\`\`bash
npm install
\`\`\`

2. Start the development server:
\`\`\`bash
npm run dev
\`\`\`

3. Open [http://localhost:5173](http://localhost:5173) in your browser

## Backend API

The application expects a backend API running at `http://localhost:8080/api` with the following endpoints:

### Authentication
- `POST /api/login` - Login with email and password
- `POST /api/register` - Register new user

### Pets
- `GET /api/pets?owner_id=0` - Get store pets
- `GET /api/pets?owner_id=me` - Get user's pets
- `POST /api/pets/:id/buy` - Purchase a pet
- `DELETE /api/pets/:id` - Delete a pet

### Products
- `GET /api/products?owner_id=0` - Get store products
- `GET /api/products?owner_id=me` - Get user's products
- `POST /api/products/:id/buy` - Purchase a product
- `DELETE /api/products/:id` - Delete a product

### Admin
- `GET /api/stats` - Get system statistics

## Design System

### Colors
- **Primary**: Warm Amber (#f59e0b) - Trust and care
- **Accent**: Teal (#14b8a6) - Nature and vitality
- **Neutrals**: Slate Gray, Light Gray, White, Dark Slate

### Typography
- **Font Family**: Inter (all weights)
- **Hierarchy**: Clear distinction between headings and body text

### Role Differentiation
- **Customer**: Warm, card-based layouts with emotional connection
- **Manager**: Sidebar navigation with teal accents, table-based data
- **Admin**: Sidebar navigation with amber accents, analytics focus

## Project Structure

\`\`\`
src/
├── components/
│   ├── layout/
│   │   ├── CustomerNav.tsx
│   │   ├── ManagerNav.tsx
│   │   └── AdminNav.tsx
│   ├── ui/
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Input.tsx
│   │   └── Modal.tsx
│   └── ItemCard.tsx
├── pages/
│   ├── Home.tsx
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── StorePets.tsx
│   ├── StoreProducts.tsx
│   ├── MyPets.tsx
│   ├── MyProducts.tsx
│   ├── Profile.tsx
│   ├── AdminDashboard.tsx
│   └── ManagerDashboard.tsx
├── stores/
│   ├── authStore.ts
│   ├── themeStore.ts
│   └── languageStore.ts
├── types/
│   └── types.ts
├── utils/
│   ├── types.ts
│   └── translations.ts
├── lib/
│   └── utils.ts
├── App.tsx
├── main.tsx
└── index.css
\`\`\`

## License

MIT
\`\`\`
