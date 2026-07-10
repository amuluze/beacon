import { describe, expect, it } from 'vitest'
import {
    containerToForm,
    emptyContainerForm,
    serializeContainerCreate,
    serializeContainerUpdate,
} from './containerForm'
import type { Container } from '@/interface/container'

describe('container form serialization', () => {
    it('serializes create fields and keeps empty collections explicit', () => {
        const form = emptyContainerForm()
        form.containerName = 'beacon'
        form.imageName = 'beacon:v3'
        form.networkName = 'bridge'
        form.restartPolicy = 'unless-stopped'

        expect(serializeContainerCreate(form, { id: 'network-id', driver: 'bridge' })).toEqual({
            container_name: 'beacon',
            image_name: 'beacon:v3',
            network_id: 'network-id',
            network_mode: 'bridge',
            network_name: 'bridge',
            restart_policy: 'unless-stopped',
            ports: [],
            volumes: [],
            environments: [],
            labels: {},
        })
    })

    it('normalizes multiline values for an update payload', () => {
        const form = emptyContainerForm()
        form.containerName = 'api'
        form.imageName = 'api:latest'
        form.networkName = 'app-net'
        form.restartPolicy = 'on-failure'
        form.ports = '8080:80\n 8443:443 \n'
        form.volumes = '/data:/app/data\n'
        form.environments = 'TZ=Asia/Shanghai\nLOG_LEVEL=info\n'
        form.labels = 'owner=platform\ninvalid\n tier = backend '

        expect(serializeContainerUpdate('abcdef123456', form)).toEqual({
            container_id: 'abcdef123456',
            container_name: 'api',
            image_name: 'api:latest',
            network_name: 'app-net',
            restart_policy: 'on-failure',
            ports: ['8080:80', '8443:443'],
            volumes: ['/data:/app/data'],
            environments: ['TZ=Asia/Shanghai', 'LOG_LEVEL=info'],
            labels: { owner: 'platform', tier: 'backend' },
        })
    })

    it('omits runtime fields that are unavailable in monitoring data', () => {
        const container = {
            id: 'abcdef',
            name: 'api',
            image: 'api:latest',
            ip: '172.18.0.2',
            ports: '8080:80',
            state: 'running',
            uptime: '2026-07-10 10:00:00',
            cpu_percent: 0,
            memory_percent: 0,
            memory_usage: 0,
            memory_limit: 0,
        } satisfies Container

        const payload = serializeContainerUpdate(container.id, containerToForm(container))

        expect(payload.ports).toEqual(['8080:80'])
        expect(payload).not.toHaveProperty('volumes')
        expect(payload).not.toHaveProperty('environments')
        expect(payload).not.toHaveProperty('labels')
    })
})
