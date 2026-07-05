/**
 * @Author     : Amu
 * @Date       : 2024/09/14 10:31:14
 * @Description:
 */

import request from '@/api'
import type { ResourceQueryResult, RoleCreateArgs, RoleDeleteArgs, RoleQueryResult, RoleUpdateArgs, UserCreateArgs, UserDeleteArgs, UserQueryResult, UserUpdateArgs } from '@/interface/account'

export async function queryUser() {
    return request.get<UserQueryResult>('/api/v1/user/user_query', {})
}

export async function createUser(params: UserCreateArgs) {
    return request.post('/api/v1/user/user_create', params)
}

export async function updateUser(params: UserUpdateArgs) {
    return request.post('/api/v1/user/user_update', params)
}

export async function deleteUser(params: UserDeleteArgs) {
    return request.post('/api/v1/user/user_delete', params)
}

export async function queryRole() {
    return request.get<RoleQueryResult>('/api/v1/role/role_query', {})
}

export async function createRole(params: RoleCreateArgs) {
    return request.post('/api/v1/role/role_create', params)
}

export async function updateRole(params: RoleUpdateArgs) {
    return request.post('/api/v1/role/role_update', params)
}

export async function deleteRole(params: RoleDeleteArgs) {
    return request.post('/api/v1/role/role_delete', params)
}

export async function queryResource() {
    return request.get<ResourceQueryResult>('/api/v1/resource/resource_query', {})
}
