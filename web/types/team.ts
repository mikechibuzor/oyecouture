// Roles the API recognises. Owner is the super admin: team, settings, delivery
// and prices. Staff runs orders and fulfilment and nothing else.
export enum AdminRole {
  Owner = 'owner',
  Staff = 'staff',
}

export interface TeamMember {
  readonly id: string
  readonly email: string
  readonly role: AdminRole
  readonly active: boolean
  readonly createdAt: string
  readonly updatedAt: string
}

export interface TeamMemberDraft {
  readonly email: string
  readonly password: string
  readonly role: AdminRole
}

export interface TeamMemberUpdate {
  readonly role: AdminRole
  readonly active: boolean
}

export interface PasswordChange {
  readonly currentPassword: string
  readonly newPassword: string
}
