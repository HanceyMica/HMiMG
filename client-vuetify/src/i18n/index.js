/**
 * 国际化 (i18n) 配置文件
 *
 * 本文件负责配置 vue-i18n 插件，实现应用的多语言支持。
 *
 * 支持的语言：
 * - en (English) - 英语
 * - zh (中文) - 简体中文
 * - ja (日本語) - 日语
 *
 * 语言文件结构：
 * - common: 通用词汇，如导航、按钮等公共文案
 * - login: 登录/注册相关文案
 * - home: 首页相关文案
 * - image: 图片详情页相关文案
 * - about: 关于页面相关文案
 * - admin: 后台管理相关文案
 * - user: 用户相关文案
 * - notFound: 404 页面文案
 */

import { createI18n } from 'vue-i18n'

/**
 * 多语言消息对象
 *
 * 每个语言对象包含多个命名空间 (namespace)，
 * 在模板中使用 $t('namespace.key') 访问对应文案
 */
const messages = {
  // 英语消息
  en: {
    common: {
      home: 'Home',
      library: 'Library',
      album: 'Album',
      collection: 'Collection',
      albums: 'Albums',
      collections: 'Collections',
      admin: 'Admin',
      about: 'About',
      logout: 'Logout',
      dark: 'Dark',
      light: 'Light',
      copyright: '{title} © 2026 HanceyMica',
      language: 'Language',
      photos: 'Images',
      edit: 'Edit',
      notLoggedIn: 'Not logged in, redirecting to login page...'
    },
    login: {
      title: 'Login',
      username: 'Username',
      password: 'Password',
      email: 'Email',
      phone: 'Phone',
      loginBtn: 'Log in',
      registerBtn: 'Register',
      toRegister: 'No account? Register now',
      toLogin: 'Already have an account? Log in',
      success: 'Login successful',
      registerSuccess: 'Registration successful',
      failed: 'Login failed',
      registerFailed: 'Registration failed',
      required: 'Please input your '
    },
    home: {
      welcome: 'Welcome to {title}',
      myLibrary: 'My Library',
      uploadPhoto: 'Upload Photo',
      selectAlbum: 'Select Album',
      selectFile: 'Select File',
      upload: 'Upload',
      uploadSuccess: 'Upload successful',
      uploadFailed: 'Upload failed',
      dragDropText: 'Click or drag file to this area to upload',
      dragDropHint: 'Support for a single or bulk upload. Strictly prohibit from uploading prohibited files!',
      selectAlbumToUpload: 'Select Album to Upload',
      noCover: 'No Cover'
    },
    image: {
      details: 'Image Details',
      filename: 'Filename',
      originalName: 'Original Name',
      size: 'Size',
      type: 'Type',
      uploadedAt: 'Uploaded At',
      album: 'Album',
      download: 'Download',
      back: 'Back',
      generateLink: 'Generate Link',
      linkSettings: 'Link Settings',
      linkType: 'Link Format',
      markdown: 'Markdown',
      html: 'HTML',
      bbcode: 'BBCode',
      direct: 'Direct Link',
      copy: 'Copy',
      copied: 'Copied!',
      edit: 'Edit',
      delete: 'Delete',
      confirmDelete: 'Are you sure you want to delete this?',
      deleteSuccess: 'Deleted successfully',
      deleteFailed: 'Failed to delete',
      previous: 'Previous',
      next: 'Next',
      nameRequired: 'Image name is required',
      update: 'Update',
      updateSuccess: 'Updated successfully',
      updateFailed: 'Failed to update'
    },
    about: {
      title: 'About HMiMG',
      description: 'HMiMG (HanceyMica Image Management Gallery) is a modern, responsive, and feature-rich image management system designed for personal and small team use.',
      version: 'Version',
      features: 'Key Features',
      feature1: 'Organize images into Albums and nested Collections',
      feature2: 'Multi-file drag-and-drop upload',
      feature3: 'Responsive design for mobile and desktop',
      feature4: 'Secure authentication and admin controls',
      author: 'Created by',
      github: 'GitHub Repository'
    },
    admin: {
      dashboard: 'Admin Dashboard',
      createAlbum: 'Create Album',
      createCollection: 'Create Collection',
      organize: 'Organize Collection',
      name: 'Name',
      description: 'Description',
      create: 'Create',
      targetCollection: 'Target Collection',
      itemType: 'Item Type',
      itemId: 'Item ID',
      itemName: 'Item Name',
      itemIdHelp: 'Enter Album ID or Collection ID',
      itemNameHelp: 'Select Album or Collection Name',
      add: 'Add',
      album: 'Album',
      collection: 'Collection',
      albumCreated: 'Album created',
      collectionCreated: 'Collection created',
      addedSuccess: 'Added to collection',
      failedCreateAlbum: 'Failed to create album',
      failedCreateCollection: 'Failed to create collection',
      failedAdd: 'Failed to add',
      systemSettings: 'System Settings',
      maxUsers: 'Max Users',
      allowRegistration: 'Allow Registration',
      yes: 'Yes',
      no: 'No',
      saveSettings: 'Save Settings',
      settingsUpdated: 'Settings updated',
      settingsFailed: 'Failed to update settings',
      websiteTitle: 'Website Title',
      websiteTitlePlaceholder: 'Enter website title (e.g. HMiMG)',
      defaultLanguage: 'Default Language',
      accountSettings: 'Account Settings',
      updateProfile: 'Update Profile',
      profileUpdated: 'Profile updated successfully',
      profileFailed: 'Failed to update profile',
      accessDenied: 'Access denied',
      passwordHelp: 'Leave blank to keep current password',
      oldPassword: 'Old Password',
      newPassword: 'New Password',
      confirmPassword: 'Confirm New Password',
      passwordMismatch: 'The new passwords do not match!',
      changePasswordSuccess: 'Password changed successfully, please login again.'
    },
    user: {
      profileUpdated: 'Profile updated successfully',
      passwordMismatch: 'The new passwords do not match!',
      changePasswordSuccess: 'Password changed successfully, please login again.'
    },
    notFound: {
      title: '404',
      description: 'Sorry, the page you visited does not exist.',
      backHome: 'Back Home'
    }
  },

  // 中文（简体）消息
  zh: {
    common: {
      home: '首页',
      library: '图库',
      album: '相册',
      collection: '合集',
      albums: '相册',
      collections: '合集',
      admin: '后台',
      about: '关于',
      logout: '登出',
      dark: '深色',
      light: '浅色',
      copyright: '{title} © 2026 HanceyMica',
      language: '语言',
      photos: '图片',
      edit: '编辑',
      notLoggedIn: '未登录，正在跳转至登录页面...'
    },
    login: {
      title: '登录',
      username: '用户名',
      password: '密码',
      email: '邮箱',
      phone: '手机号',
      loginBtn: '登录',
      registerBtn: '注册',
      toRegister: '没有账号？立即注册',
      toLogin: '已有账号？立即登录',
      success: '登录成功',
      registerSuccess: '注册成功',
      failed: '登录失败',
      registerFailed: '注册失败',
      required: '请输入您的'
    },
    home: {
      welcome: '欢迎使用 {title}',
      myLibrary: '我的图库',
      uploadPhoto: '上传照片',
      selectAlbum: '选择相册',
      selectFile: '选择文件',
      upload: '上传',
      uploadSuccess: '上传成功',
      uploadFailed: '上传失败',
      dragDropText: '点击或拖拽文件到此区域上传',
      dragDropHint: '支持单次或批量上传。严禁上传违禁文件！',
      selectAlbumToUpload: '选择要上传的相册',
      noCover: '暂无封面'
    },
    image: {
      details: '图片详情',
      filename: '文件名',
      originalName: '原文件名',
      size: '大小',
      type: '类型',
      uploadedAt: '上传时间',
      album: '所属相册',
      download: '下载',
      back: '返回',
      generateLink: '生成外链',
      linkSettings: '外链设置',
      linkType: '链接格式',
      markdown: 'Markdown',
      html: 'HTML',
      bbcode: 'BBCode',
      direct: '直链',
      copy: '复制',
      copied: '已复制!',
      edit: '编辑',
      delete: '删除',
      confirmDelete: '确定要删除吗？',
      deleteSuccess: '删除成功',
      deleteFailed: '删除失败',
      previous: '上一张',
      next: '下一张',
      nameRequired: '请输入图片名称',
      update: '更新',
      updateSuccess: '更新成功',
      updateFailed: '更新失败'
    },
    about: {
      title: '关于 HMiMG',
      description: 'HMiMG (HanceyMica Image Management Gallery) 是一个专为个人和小型团队设计的现代、响应式且功能丰富的图片管理系统。',
      version: '版本',
      features: '主要功能',
      feature1: '将图片整理到相册和嵌套合集中',
      feature2: '多文件拖拽上传',
      feature3: '适配移动端和桌面的响应式设计',
      feature4: '安全的身份验证和后台管理',
      author: '作者',
      github: 'GitHub 仓库'
    },
    admin: {
      dashboard: '后台管理',
      createAlbum: '创建相册',
      createCollection: '创建合集',
      organize: '管理合集',
      name: '名称',
      description: '描述',
      create: '创建',
      targetCollection: '目标合集',
      itemType: '项目类型',
      itemId: '项目ID',
      itemName: '项目名称',
      itemIdHelp: '输入相册ID或合集ID',
      itemNameHelp: '选择相册或合集名称',
      add: '添加',
      album: '相册',
      collection: '合集',
      albumCreated: '相册已创建',
      collectionCreated: '合集已创建',
      addedSuccess: '已添加到合集',
      failedCreateAlbum: '创建相册失败',
      failedCreateCollection: '创建合集失败',
      failedAdd: '添加失败',
      systemSettings: '系统设置',
      maxUsers: '最大用户数',
      allowRegistration: '允许注册',
      yes: '是',
      no: '否',
      saveSettings: '保存设置',
      settingsUpdated: '设置已更新',
      settingsFailed: '更新设置失败',
      websiteTitle: '网站标题',
      websiteTitlePlaceholder: '输入网站标题 (例如 HMiMG)',
      defaultLanguage: '默认界面语言',
      accountSettings: '账户设置',
      updateProfile: '更新个人信息',
      profileUpdated: '个人信息已更新',
      profileFailed: '更新个人信息失败',
      accessDenied: '访问被拒绝',
      passwordHelp: '留空以保留当前密码',
      oldPassword: '旧密码',
      newPassword: '新密码',
      confirmPassword: '确认新密码',
      passwordMismatch: '两次输入的新密码不一致！',
      changePasswordSuccess: '密码修改成功，请重新登录。'
    },
    user: {
      profileUpdated: '个人信息已更新',
      passwordMismatch: '两次输入的新密码不一致！',
      changePasswordSuccess: '密码修改成功，请重新登录。'
    },
    notFound: {
      title: '404',
      description: '抱歉，您访问的页面不存在。',
      backHome: '返回首页'
    }
  },

  // 日语消息
  ja: {
    common: {
      home: 'ホーム',
      library: 'ライブラリ',
      album: 'アルバム',
      collection: 'コレクション',
      albums: 'アルバム',
      collections: 'コレクション',
      admin: '管理',
      about: '情報',
      logout: 'ログアウト',
      dark: 'ダーク',
      light: 'ライト',
      copyright: '{title} © 2026 HanceyMica',
      language: '言語',
      photos: '画像',
      edit: '編集',
      notLoggedIn: 'ログインしていません。ログインページにリダイレクトしています...'
    },
    login: {
      title: 'ログイン',
      username: 'ユーザー名',
      password: 'パスワード',
      email: 'メールアドレス',
      phone: '電話番号',
      loginBtn: 'ログイン',
      registerBtn: '登録',
      toRegister: 'アカウントをお持ちでないですか？登録する',
      toLogin: 'すでにアカウントをお持ちですか？ログインする',
      success: 'ログイン成功',
      registerSuccess: '登録成功',
      failed: 'ログイン失敗',
      registerFailed: '登録失敗',
      required: '入力してください：'
    },
    home: {
      welcome: '{title} へようこそ',
      myLibrary: 'マイライブラリ',
      uploadPhoto: '写真をアップロード',
      selectAlbum: 'アルバムを選択',
      selectFile: 'ファイルを選択',
      upload: 'アップロード',
      uploadSuccess: 'アップロード成功',
      uploadFailed: 'アップロード失敗',
      dragDropText: 'クリックまたはドラッグ＆ドロップでアップロード',
      dragDropHint: '単一または複数ファイルのアップロードをサポートしています。違法ファイルのアップロードは固くお断りします！',
      selectAlbumToUpload: 'アップロードするアルバムを選択',
      noCover: 'カバーなし'
    },
    image: {
      details: '画像詳細',
      filename: 'ファイル名',
      originalName: '元のファイル名',
      size: 'サイズ',
      type: 'タイプ',
      uploadedAt: 'アップロード日時',
      album: 'アルバム',
      download: 'ダウンロード',
      back: '戻る',
      generateLink: 'リンク生成',
      linkSettings: 'リンク設定',
      linkType: 'リンク形式',
      markdown: 'Markdown',
      html: 'HTML',
      bbcode: 'BBCode',
      direct: '直リンク',
      copy: 'コピー',
      copied: 'コピーしました！',
      edit: '編集',
      delete: '削除',
      confirmDelete: '本当に削除しますか？',
      deleteSuccess: '削除成功',
      deleteFailed: '削除失敗',
      previous: '前の画像',
      next: '次の画像',
      nameRequired: '画像名を入力してください',
      update: '更新',
      updateSuccess: '更新成功',
      updateFailed: '更新失敗'
    },
    about: {
      title: 'HMiMG について',
      description: 'HMiMG (HanceyMica Image Management Gallery) は、個人や小規模チーム向けに設計された、モダンでレスポンシブ、かつ機能豊富な画像管理システムです。',
      version: 'バージョン',
      features: '主な機能',
      feature1: 'アルバムとネストされたコレクションへの画像整理',
      feature2: '複数ファイルのドラッグ＆ドロップアップロード',
      feature3: 'モバイルとデスクトップに対応したレスポンシブデザイン',
      feature4: '安全な認証と管理機能',
      author: '作成者',
      github: 'GitHub リポジトリ'
    },
    admin: {
      dashboard: '管理ダッシュボード',
      createAlbum: 'アルバム作成',
      createCollection: 'コレクション作成',
      organize: 'コレクション整理',
      name: '名前',
      description: '説明',
      create: '作成',
      targetCollection: '対象コレクション',
      itemType: 'アイテムタイプ',
      itemId: 'アイテムID',
      itemName: 'アイテム名',
      itemIdHelp: 'アルバムIDまたはコレクションIDを入力',
      itemNameHelp: 'アルバムまたはコレクション名を選択',
      add: '追加',
      album: 'アルバム',
      collection: 'コレクション',
      albumCreated: 'アルバムが作成されました',
      collectionCreated: 'コレクションが作成されました',
      addedSuccess: 'コレクションに追加されました',
      failedCreateAlbum: 'アルバムの作成に失敗しました',
      failedCreateCollection: 'コレクションの作成に失敗しました',
      failedAdd: '追加に失敗しました',
      systemSettings: 'システム設定',
      maxUsers: '最大ユーザー数',
      allowRegistration: '登録を許可',
      yes: 'はい',
      no: 'いいえ',
      saveSettings: '設定を保存',
      settingsUpdated: '設定が更新されました',
      settingsFailed: '設定の更新に失敗しました',
      websiteTitle: 'ウェブサイトのタイトル',
      websiteTitlePlaceholder: 'ウェブサイトのタイトルを入力 (例: HMiMG)',
      defaultLanguage: 'デフォルトのインターフェース言語',
      accountSettings: 'アカウント設定',
      updateProfile: 'プロフィール更新',
      profileUpdated: 'プロフィールが更新されました',
      profileFailed: 'プロフィールの更新に失敗しました',
      accessDenied: 'アクセスが拒否されました',
      passwordHelp: '現在のパスワードを保持する場合は空白のままにしてください',
      oldPassword: '現在のパスワード',
      newPassword: '新しいパスワード',
      confirmPassword: '新しいパスワード（確認）',
      passwordMismatch: '新しいパスワードが一致しません！',
      changePasswordSuccess: 'パスワードが変更されました。再度ログインしてください。'
    },
    user: {
      profileUpdated: 'プロフィールが更新されました',
      passwordMismatch: '新しいパスワードが一致しません！',
      changePasswordSuccess: 'パスワードが変更されました。再度ログインしてください。'
    },
    notFound: {
      title: '404',
      description: '申し訳ありませんが、アクセスしたページは存在しません。',
      backHome: 'ホームに戻る'
    }
  }
}

/**
 * 获取初始语言设置
 *
 * 优先级顺序：
 * 1. localStorage 中存储的用户语言偏好 (lang)
 * 2. 默认使用中文 (zh)
 *
 * @returns {string} 语言代码，如 'zh'、'en' 或 'ja'
 */
const getInitialLocale = () => {
  if (typeof window !== 'undefined') {
    // 尝试从本地存储获取用户之前设置的语言
    return localStorage.getItem('lang') || 'zh'
  }
  // 服务端渲染时默认返回中文
  return 'zh'
}

/**
 * 创建 i18n 实例
 *
 * 配置选项说明：
 * - legacy: false - 使用 Composition API 模式（推荐）
 * - locale: 初始语言，从 getInitialLocale() 获取
 * - fallbackLocale: 当某个翻译 key 不存在时使用的备用语言
 * - messages: 包含所有语言翻译的对象
 */
const i18n = createI18n({
  legacy: false,
  locale: getInitialLocale(),
  fallbackLocale: 'en',
  messages,
})

// 导出 i18n 实例，供 main.js 中注册到 Vue 应用
export default i18n
