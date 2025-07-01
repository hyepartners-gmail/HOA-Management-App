# Bear Paw HOA Management App

A full-featured HOA management system built with:

- **Frontend**: React + TypeScript (Vite, TailwindCSS, Radix UI, Lucide Icons)
- **Backend**: Go (REST API using chi)
- **Database**: Google Cloud Datastore
- **Auth**: JWT-based
- **Email**: SendGrid or Mailgun
- **SMS**: Twilio
- **PDF Generation**: gofpdf + Google Cloud Storage
- **Hosting**: Google Cloud Run (single container for frontend + backend)

---

## 🏗️ Features

- Role-based dashboard: Admin, President, Treasurer, Secretary, Cabin Owner
- Auth (login, reset password, RBAC)
- Cabin and Owner management
- Invoice generation + Stripe payments
- Meeting minutes, proxy voting, service requests
- Flash notifications via SMS/email
- Document center, audit logs, community polls
- Public site + talent directory + message board

---

## 🚀 Quickstart (Local Dev)

### Prerequisites

- Go 1.22+
- Node 20+
- Docker (for Cloud Run emulation or building locally)
- Firebase/Datastore emulator (optional for full local backend testing)

### 1. Clone the Repo

```bash
git clone git@github.com:your-org/HOA-Management-App.git
cd HOA-Management-App
````

### 2. Install Frontend Dependencies

```bash
cd frontend
npm install
npm run dev
```

### 3. Start Go Backend

```bash
cd backend
go mod tidy
go run main.go
```

* The Go server serves API routes under `/api/*`
* It also statically serves the React app under `/`

### 4. Set Environment Variables (Local Dev)

Create a `.env` file in the root:

```env
JWT_SECRET=super-secret-key
SENDGRID_KEY=your-sendgrid-key
MAILGUN_API_KEY=your-mailgun-key
STRIPE_SECRET=your-stripe-secret
TWILIO_SID=...
TWILIO_TOKEN=...
TWILIO_FROM_NUMBER=...
GCS_BUCKET_NAME=...
```

---

## 🐳 Deployment on Google Cloud Run

This project is designed to be deployed as a **single container** on Cloud Run.

### Dockerfile already included

Build + deploy:

```bash
gcloud run deploy bearpaw-app \
  --source . \
  --region us-west1 \
  --allow-unauthenticated
```

---

## 🧪 Testing

Backend tests (if written):

```bash
go test ./...
```

---

## 📁 File Structure

```
HOA-Management-App/
├── frontend/        # Vite React app (UI)
├── backend/         # Go source code
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── utils/
├── Dockerfile
├── go.mod
├── .env
└── README.md
```

---

## 👥 Contributing

1. Fork the repo
2. Create a new branch: `git checkout -b feature-name`
3. Commit your changes
4. Push and open a PR

---

## 📄 License

MIT 

