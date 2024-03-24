<template>
    <el-card shadow="never">
        <el-table :data="imageData" border stripe ref="multipleTable" v-loading="loading">
            <el-table-column prop="id" label="镜像 ID" align="center" fixed />
            <el-table-column prop="name" label="镜像名称" align="center" min-width="100" show-overflow-tooltip fixed />
            <el-table-column prop="tag" label="镜像Tag" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="created" label="创建时间" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="size" label="镜像大小" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column label="操作" width="180" fixed="right" align="center">
                <template #default="scope">
                    <el-button type="primary" size="default" @click="deleteImage(scope.row.id)">
                        <i-ep-view />删除
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-card>
</template>

<script setup lang="ts"></script>

<style scoped lang="scss"></style>
