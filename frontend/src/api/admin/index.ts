/**
 * Admin API barrel export
 * Centralized exports for all admin API modules
 */

import dashboardAPI from './dashboard'
import usersAPI from './users'
import announcementsAPI from './announcements'
import settingsAPI from './settings'
import systemAPI from './system'
import userAttributesAPI from './userAttributes'
import backupAPI from './backup'

/**
 * Unified admin API object for convenient access
 */
export const adminAPI = {
  dashboard: dashboardAPI,
  users: usersAPI,
  announcements: announcementsAPI,
  settings: settingsAPI,
  system: systemAPI,
  userAttributes: userAttributesAPI,
  backup: backupAPI
}

export {
  dashboardAPI,
  usersAPI,
  announcementsAPI,
  settingsAPI,
  systemAPI,
  userAttributesAPI,
  backupAPI
}

export default adminAPI

// Re-export types used by components
export type { BalanceHistoryItem } from './users'
