export default {
  // 路由国际化
  route: {
    dashboard: "首页",
    document: "项目文档",
    userset: "修改密码",
    system:"系统管理",
    nodelist:"节点列表",
    sublist:"订阅列表",
    subcription:"订阅管理",
    templatelist:"模板列表",
    apikey: "API密钥管理",
  },
  // 登录页面国际化
  login: {
    username: "用户名",
    password: "密码",
    login: "登 录",
    captchaCode: "验证码",
    capsLock: "大写锁定已打开",
    message: {
      username: {
        required: "请输入用户名",
      },
      password: {
        required: "请输入密码",
        min: "密码不能少于6位",
      },
      captchaCode: {
        required: "请输入验证码",
      },
    },
  },
    // 重置密码页面国际化
    userset:{
      title: "修改密码",
      newUsername: "新账号",
      newPassword: "新密码",
      message: {
        title:"提示",
        xx1:"账号或密码不能为空",
        xx2: "密码长度不能小于6位",
        xx3:"你确定要重置密码吗",
        xx4:"密码重置成功，新密码是：",
      },
    },
  // 导航栏国际化
  navbar: {
    dashboard: "首页",
    logout: "注销登出",
    userset: "修改密码",
  },
  sizeSelect: {
    tooltip: "布局大小",
    default: "默认",
    large: "大型",
    small: "小型",
    message: {
      success: "切换布局大小成功！",
    },
  },
  langSelect: {
    message: {
      success: "切换语言成功！",
    },
  },  settings: {
    project: "项目配置",
    theme: "主题设置",
    interface: "界面设置",
    navigation: "导航设置",
    themeColor: "主题颜色",
    tagsView: "开启 Tags-View",
    fixedHeader: "固定 Header",
    sidebarLogo: "侧边栏 Logo",
    watermark: "开启水印",
  },
  // 通用按钮文本
  confirm: "确定",
  cancel: "取消",
  // API密钥管理页面国际化
  apikey: {
    title: "API密钥管理",
    createNew: "创建新密钥",
    description: "描述",
    descriptionPlaceholder: "输入密钥用途描述",
    createdAt: "创建时间",
    expiredAt: "过期时间",
    actions: "操作",
    delete: "删除",
    noData: "暂无数据",
    expiration: "过期时间",
    neverExpire: "永不过期",
    customExpire: "自定义",
    selectExpiration: "选择过期时间",
    newKeyCreated: "新密钥已创建",
    saveKeyWarning: "请保存此密钥，它只会显示一次！",
    copy: "复制",
    iSavedIt: "我已保存",
    fetchFailed: "获取API密钥列表失败",
    createFailed: "创建API密钥失败",
    deleteFailed: "删除API密钥失败",
    deleteSuccess: "API密钥删除成功",
    copySuccess: "复制成功",
    search: "搜索密钥",
    deleteConfirmTitle: "删除确认",
    deleteConfirmMessage: "确定要删除此API密钥吗？此操作不可撤销。",
    deleting: "正在删除...",
    descriptionRequired: "请输入描述",
    keyCreatedButNotRetrieved: "密钥创建成功，但无法获取密钥值",
    copyFailed: "复制失败，请手动复制",
    manage: "管理API密钥"
  },
}
