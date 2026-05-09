import { BriefcaseBusiness, LogOut, Users } from "lucide-react"
import { NavLink, Outlet, useNavigate } from "react-router-dom"

import { Button } from "@/components/ui/button"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { logout } from "@/lib/auth"

const navigationItems = [
  { to: "/campaigns", label: "Campaigns", icon: BriefcaseBusiness },
  { to: "/contacts", label: "Contacts", icon: Users },
]

export function AppLayout() {
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate("/login", { replace: true })
  }

  return (
    <div className="flex min-h-screen bg-background">
      <Sidebar>
        <SidebarHeader>
          <div className="space-y-1">
            <p className="text-lg font-semibold">NudgeBuddy</p>
            <p className="text-muted-foreground text-sm">Outreach control center</p>
          </div>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Workspace</SidebarGroupLabel>
            <SidebarMenu>
              {navigationItems.map((item) => {
                const Icon = item.icon

                return (
                  <SidebarMenuItem key={item.to}>
                    <NavLink to={item.to}>
                      {({ isActive }) => (
                        <SidebarMenuButton isActive={isActive}>
                          <Icon className="size-4" />
                          <span>{item.label}</span>
                        </SidebarMenuButton>
                      )}
                    </NavLink>
                  </SidebarMenuItem>
                )
              })}
            </SidebarMenu>
          </SidebarGroup>
          <div className="mt-auto px-2">
            <Button className="w-full justify-start" variant="ghost" onClick={handleLogout}>
              <LogOut className="size-4" />
              Logout
            </Button>
          </div>
        </SidebarContent>
      </Sidebar>
      <SidebarInset className="p-6 md:p-8">
        <Outlet />
      </SidebarInset>
    </div>
  )
}
