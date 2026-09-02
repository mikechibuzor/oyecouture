import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  PasswordChange,
  TeamMember,
  TeamMemberDraft,
  TeamMemberUpdate,
} from '~/types/team'

export const useAdminTeam = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchTeam = async (): Promise<readonly TeamMember[]> => {
    const response = await $fetch<{ members: TeamMember[] }>(ADMIN_API_ROUTES.team, {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

    return response.members
  }

  const createMember = async (draft: TeamMemberDraft): Promise<TeamMember> => {
    const response = await $fetch<{ member: TeamMember }>(ADMIN_API_ROUTES.team, {
      baseURL: apiBaseUrl,
      body: draft,
      credentials: 'include',
      method: 'POST',
    })

    return response.member
  }

  const updateMember = async (
    id: string,
    update: TeamMemberUpdate,
  ): Promise<TeamMember> => {
    const response = await $fetch<{ member: TeamMember }>(ADMIN_API_ROUTES.teamMember(id), {
      baseURL: apiBaseUrl,
      body: update,
      credentials: 'include',
      method: 'PATCH',
    })

    return response.member
  }

  const deleteMember = (id: string): Promise<void> =>
    $fetch(ADMIN_API_ROUTES.teamMember(id), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'DELETE',
    })

  const changePassword = (change: PasswordChange): Promise<void> =>
    $fetch(ADMIN_API_ROUTES.password, {
      baseURL: apiBaseUrl,
      body: change,
      credentials: 'include',
      method: 'PATCH',
    })

  return { fetchTeam, createMember, updateMember, deleteMember, changePassword }
}
