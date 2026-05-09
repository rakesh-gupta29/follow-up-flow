import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"

import { AuthGuard } from "./auth-guard"
import { AppLayout } from "./layouts/app-layout"
import { CampaignContactsPage } from "./pages/campaign-contacts-page"
import { CampaignsPage } from "./pages/campaigns-page"
import { ContactsPage } from "./pages/contacts-page"
import { LoginPage } from "./pages/login-page"
import { ProfilePage } from "./pages/profile-page"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AuthGuard><AppLayout /></AuthGuard>}>
          <Route path="/" element={<Navigate to="/campaigns" replace />} />
          <Route path="/campaigns" element={<CampaignsPage />} />
          <Route path="/campaigns/:campaignId/contacts" element={<CampaignContactsPage />} />
          <Route path="/contacts" element={<ContactsPage />} />
          <Route path="/profile" element={<ProfilePage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
