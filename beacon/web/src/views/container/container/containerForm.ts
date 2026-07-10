import type {
    Container,
    CreateContainerArgs,
    RestartPolicy,
    UpdateContainerArgs,
} from '@/interface/container'

export interface ContainerFormModel {
    containerName: string
    imageName: string
    networkName: string
    restartPolicy: RestartPolicy
    ports: string
    volumes: string
    environments: string
    labels: string
    preserveUnavailable: {
        volumes: boolean
        environments: boolean
        labels: boolean
    }
}

interface NetworkSelection {
    id: string
    driver: string
}

export const restartPolicyOptions: RestartPolicy[] = [
    'no',
    'always',
    'unless-stopped',
    'on-failure',
]

export function emptyContainerForm(): ContainerFormModel {
    return {
        containerName: '',
        imageName: '',
        networkName: '',
        restartPolicy: 'always',
        ports: '',
        volumes: '',
        environments: '',
        labels: '',
        preserveUnavailable: {
            volumes: false,
            environments: false,
            labels: false,
        },
    }
}

function splitLines(value: string): string[] {
    return value
        .split(/\r?\n/)
        .map(item => item.trim())
        .filter(Boolean)
}

function parseLabels(value: string): Record<string, string> {
    return Object.fromEntries(splitLines(value).flatMap((line) => {
        const separator = line.indexOf('=')
        if (separator <= 0)
            return []
        const key = line.slice(0, separator).trim()
        const labelValue = line.slice(separator + 1).trim()
        return key ? [[key, labelValue]] : []
    }))
}

function serializeShared(form: ContainerFormModel) {
    return {
        container_name: form.containerName.trim(),
        image_name: form.imageName.trim(),
        network_name: form.networkName.trim(),
        restart_policy: form.restartPolicy,
        ports: splitLines(form.ports),
        volumes: splitLines(form.volumes),
        environments: splitLines(form.environments),
        labels: parseLabels(form.labels),
    }
}

export function serializeContainerCreate(form: ContainerFormModel, network?: NetworkSelection): CreateContainerArgs {
    return {
        ...serializeShared(form),
        network_id: network?.id ?? '',
        network_mode: network?.driver ?? 'bridge',
    }
}

export function serializeContainerUpdate(containerID: string, form: ContainerFormModel): UpdateContainerArgs {
    const shared = serializeShared(form)
    return {
        container_id: containerID,
        container_name: shared.container_name,
        image_name: shared.image_name,
        network_name: shared.network_name,
        restart_policy: shared.restart_policy,
        ports: shared.ports,
        ...form.preserveUnavailable.volumes ? {} : { volumes: shared.volumes },
        ...form.preserveUnavailable.environments ? {} : { environments: shared.environments },
        ...form.preserveUnavailable.labels ? {} : { labels: shared.labels },
    }
}

export function containerToForm(container: Container): ContainerFormModel {
    const ports = container.ports
        .split(',')
        .map(item => item.trim())
        .filter(Boolean)
        .join('\n')
    const labels = Object.entries(container.labels ?? {})
        .map(([key, value]) => `${key}=${value}`)
        .join('\n')

    const form = {
        ...emptyContainerForm(),
        containerName: container.name,
        imageName: container.image,
        networkName: container.network ?? '',
        ports,
        volumes: container.volumes ?? '',
        environments: container.environments ?? '',
        labels,
    }
    form.preserveUnavailable.volumes = container.volumes === undefined
    form.preserveUnavailable.environments = container.environments === undefined
    form.preserveUnavailable.labels = container.labels === undefined
    return form
}
