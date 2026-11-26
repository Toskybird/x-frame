<template>
  <AppLayout>
    <TablePageLayout>
      <template #actions>
        <div class="flex flex-wrap items-center justify-end gap-3">
          <button class="btn btn-secondary" @click="showAttributesModal = true">
            {{ t('admin.users.attributes.title') }}
          </button>
          <button class="btn btn-primary" @click="showCreateModal = true">
            <Icon name="plus" size="sm" class="mr-1" />
            {{ t('admin.users.createUser') }}
          </button>
        </div>
      </template>

      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <input
            v-model="searchQuery"
            type="search"
            class="input max-w-xs"
            :placeholder="t('common.search')"
            @keyup.enter="handleSearch"
          />
          <Select
            v-model="filters.role"
            :options="roleOptions"
            class="w-36"
            @update:model-value="handleSearch"
          />
          <Select
            v-model="filters.status"
            :options="statusOptions"
            class="w-36"
            @update:model-value="handleSearch"
          />
          <button class="btn btn-secondary" @click="handleSearch">{{ t('common.search') }}</button>
        </div>
      </template>

      <template #table>
        <div v-if="loading" class="flex justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
        </div>
        <EmptyState v-else-if="users.length === 0" :title="t('common.noData')" />
        <DataTable
          v-else
          :columns="columns"
          :data="users"
          :sort-by="sortBy"
          :sort-order="sortOrder"
          @sort="handleSort"
        >
          <template #cell-email="{ row }">
            <div class="min-w-0">
              <div class="truncate font-medium text-gray-900 dark:text-white">{{ row.email }}</div>
              <div v-if="row.username" class="truncate text-xs text-gray-500">{{ row.username }}</div>
            </div>
          </template>
          <template #cell-role="{ value }">
            <span :class="['badge', value === 'admin' ? 'badge-purple' : 'badge-gray']">
              {{ t('admin.users.roles.' + value) }}
            </span>
          </template>
          <template #cell-status="{ value }">
            <div class="flex items-center gap-1.5">
              <span
                :class="[
                  'inline-block h-2 w-2 rounded-full',
                  value === 'active' ? 'bg-green-500' : 'bg-red-500'
                ]"
              ></span>
              <span class="text-sm text-gray-700 dark:text-gray-300">
                {{ value === 'active' ? t('common.active') : t('admin.users.disabled') }}
              </span>
            </div>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-600 dark:text-gray-300">{{ formatDateTime(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <button class="btn btn-secondary btn-sm" @click="handleEdit(row)">{{ t('common.edit') }}</button>
              <button
                v-if="row.role !== 'admin'"
                class="btn btn-secondary btn-sm text-red-600"
                @click="handleDelete(row)"
              >
                {{ t('common.delete') }}
              </button>
            </div>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <ConfirmDialog
      :show="showDeleteDialog"
      :title="t('admin.users.deleteUser')"
      :message="t('admin.users.deleteConfirm', { email: deletingUser?.email })"
      :danger="true"
      @confirm="confirmDelete"
      @cancel="showDeleteDialog = false"
    />
    <UserCreateModal :show="showCreateModal" @close="showCreateModal = false" @success="loadUsers" />
    <UserEditModal :show="showEditModal" :user="editingUser" @close="closeEditModal" @success="loadUsers" />
    <UserAttributesConfigModal :show="showAttributesModal" @close="handleAttributesModalClose" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import { adminAPI } from '@/api/admin'
import type { AdminUser } from '@/types'
import type { Column } from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import UserAttributesConfigModal from '@/components/user/UserAttributesConfigModal.vue'
import UserCreateModal from '@/components/admin/user/UserCreateModal.vue'
import UserEditModal from '@/components/admin/user/UserEditModal.vue'

const { t } = useI18n()
const appStore = useAppStore()

const columns = computed<Column[]>(() => [
  { key: 'email', label: t('admin.users.columns.user'), sortable: true },
  { key: 'id', label: 'ID', sortable: true },
  { key: 'role', label: t('admin.users.columns.role'), sortable: true },
  { key: 'status', label: t('admin.users.columns.status'), sortable: true },
  { key: 'created_at', label: t('admin.users.columns.created'), sortable: true },
  { key: 'actions', label: t('admin.users.columns.actions'), sortable: false }
])

const users = ref<AdminUser[]>([])
const loading = ref(false)
const searchQuery = ref('')
const sortBy = ref('created_at')
const sortOrder = ref<'asc' | 'desc'>('desc')
const pagination = reactive({ page: 1, page_size: getPersistedPageSize(), total: 0 })
const filters = reactive({ role: '', status: '' })

const roleOptions = computed(() => [
  { value: '', label: t('common.all') },
  { value: 'admin', label: t('admin.users.roles.admin') },
  { value: 'user', label: t('admin.users.roles.user') }
])
const statusOptions = computed(() => [
  { value: '', label: t('common.all') },
  { value: 'active', label: t('common.active') },
  { value: 'disabled', label: t('admin.users.disabled') }
])

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingUser = ref<AdminUser | null>(null)
const showDeleteDialog = ref(false)
const deletingUser = ref<AdminUser | null>(null)
const showAttributesModal = ref(false)

async function loadUsers() {
  loading.value = true
  try {
    const res = await adminAPI.users.list(pagination.page, pagination.page_size, {
      search: searchQuery.value || undefined,
      role: (filters.role || undefined) as 'admin' | 'user' | undefined,
      status: (filters.status || undefined) as 'active' | 'disabled' | undefined
    })
    users.value = res.items || []
    pagination.total = res.total || 0
  } catch (error: any) {
    console.error('Failed to load users:', error)
    appStore.showError(error.response?.data?.detail || t('common.error'))
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  loadUsers()
}
function handleSort(payload: { key: string } | string) {
  const key = typeof payload === 'string' ? payload : payload.key
  if (sortBy.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = key
    sortOrder.value = 'asc'
  }
  loadUsers()
}
function handlePageChange(page: number) {
  pagination.page = page
  loadUsers()
}
function handlePageSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  loadUsers()
}
function handleEdit(user: AdminUser) {
  editingUser.value = user
  showEditModal.value = true
}
function closeEditModal() {
  showEditModal.value = false
  editingUser.value = null
}
function handleDelete(user: AdminUser) {
  deletingUser.value = user
  showDeleteDialog.value = true
}
async function confirmDelete() {
  if (!deletingUser.value) return
  try {
    await adminAPI.users.delete(deletingUser.value.id)
    appStore.showSuccess(t('common.success'))
    showDeleteDialog.value = false
    deletingUser.value = null
    await loadUsers()
  } catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('common.error'))
  }
}
async function handleAttributesModalClose() {
  showAttributesModal.value = false
  await loadUsers()
}

onMounted(() => {
  loadUsers()
})
</script>
