<template>
    <el-card shadow="never">
        <el-table :data="imageData" border stripe ref="multipleTable" v-loading="loading">
            <el-table-column prop="id" label="镜像 ID" align="center" fixed />
            <el-table-column prop="name" label="镜像名称" align="center" min-width="100" show-overflow-tooltip fixed />
            <el-table-column prop="tag" label="镜像 Tag" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="created" label="创建时间" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="size" label="镜像大小" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column label="操作" width="180" fixed="right" align="center">
                <template #default="scope">
                    <el-button type="primary" size="default" @click="deleteImage(scope.row.id)">
                        <i-ep-delete />删除
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-card>
</template>

<script setup lang="ts">
import { queryImages } from '@/api/container'
import { warning } from '@/components/Message/message.ts'
import { Image } from '@/interface/container'

const loading = ref(true)

type imageDataType = {
    id: string
    name: string
    tag: string
    created: string
    size: string
}

const imageData = ref<imageDataType[]>([])

onMounted(() => {
    getData()
})

const getData = async () => {
    loading.value = true
    const { data } = await queryImages()
    data.images.map((item: Image) => {
        const imageItem: imageDataType = {
            id: item.id,
            name: item.name,
            tag: item.tag,
            created: item.created,
            size: item.size
        }
        imageData.value.push(imageItem)
    })
    loading.value = false
}

const deleteImage = (id: string) => {
    console.log(id)
    warning('该功能暂未实现')
}
</script>

<style scoped lang="scss"></style>
