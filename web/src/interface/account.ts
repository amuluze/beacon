/**
 * @Author     : Amu
 * @Date       : 2024/09/14 10:33:31
 * @Description:
 */

export interface User {
    id: string
    username: string
    remark: string
    is_admin: number
    status: number
    created_at: string
    roles: Role[]
}

export interface UserCreateArgs {
    username: string
    password: string
    remark: string
    status: number
    role_ids: string[]
}

export interface UserUpdateArgs {
    id: string
    username: string
    remark: string
    status: number
    role_ids: string[]
}

export interface UserDeleteArgs {
    ids: string[]
}

export interface UserQueryResult {
    data: User[]
    total: number
    page: number
    size: number
}

export interface Role {
    id: string
    name: string
    remark: string
    status: number
    created_at: string
    resources: Resource[]
}

export interface RoleCreateArgs {
    name: string
    remark: string
    status: number
    resource_ids: string[]
}

export interface RoleUpdateArgs {
    id: string
    name: string
    remark: string
    status: number
    resource_ids: string[]
}

export interface RoleDeleteArgs {
    ids: string[]
}

export interface RoleQueryResult {
    data: Role[]
    total: number
    page: number
    size: number
}

export interface Resource {
    id: string
    name: string
    path: string
    method: string
    status: number
}

export interface ResourceQueryResult {
    data: Resource[]
    total: number
    page: number
    size: number
}
