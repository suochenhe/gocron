<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main>
      <notification-tab></notification-tab>
      <el-form ref="form" :model="form" :rules="formRules" label-width="180px" style="width: 700px;">
        <el-form-item label="DingTalk URL" prop="url">
          <el-input v-model="form.url"></el-input>
        </el-form-item>
        <el-form-item label="模板" prop="template">
          <el-input
            type="textarea"
            :rows="8"
            placeholder=""
            size="medium"
            v-model="form.template">
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit">保存</el-button>
        </el-form-item>
        <h3>通知用户</h3>
        <el-button type="primary" @click="createUser">新增用户</el-button> <br><br>
        <el-tag
          v-for="item in users"
          :key="item.id"
          closable
          @close="deleteUser(item)"
        >
          {{item.name}}
        </el-tag>
      </el-form>
      <el-dialog
        title=""
        :visible.sync="dialogVisible"
        width="30%">
        <el-form :model="form">
          <el-form-item label="姓名" >
            <el-input v-model.trim="name" v-focus></el-input>
          </el-form-item>
          <el-form-item label="电话" >
            <el-input v-model.trim="mobile" v-focus></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveUser">确 定</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import systemSidebar from '../sidebar'
import notificationTab from './tab'
import notificationService from '../../../api/notification'
export default {
  name: 'notification-ding',
  data () {
    return {
      dialogVisible: false,
      form: {
        url: '',
        template: ''
      },
      formRules: {
        url: [
          {type: 'url', required: true, message: '请输入有效的URL', trigger: 'blur'}
        ],
        template: [
          {required: true, message: '请输入通知模板', trigger: 'blur'}
        ]
      },
      users: [],
      name: '',
      mobile: ''
    }
  },
  components: {notificationTab, systemSidebar},
  created () {
    this.init()
  },
  methods: {
    createUser () {
      this.dialogVisible = true
    },
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      notificationService.updateDing(this.form, () => {
        this.$message.success('更新成功')
        this.init()
      })
    },
    isMobile (value) {
      return /^1[3-9]\d{9}$/.test(value)
    },
    saveUser () {
      if (this.name === '') {
        this.$message.error('请输入姓名')
        return
      }
      if (this.mobile === '') {
        this.$message.error('请输入手机号')
        return
      } else {
        if (!this.isMobile(this.mobile)) {
          this.$message.error('请输入正确的手机号')
          return
        }
      }
      notificationService.createDingUser(this.name, this.mobile, () => {
        this.dialogVisible = false
        this.init()
      })
    },
    deleteUser (item) {
      notificationService.removeDingUser(item.id, () => {
        this.init()
      })
    },
    init () {
      this.name = ''
      this.mobile = ''
      notificationService.ding((data) => {
        this.form.url = data.url
        this.form.template = data.template
        this.users = data.users
      })
    }
  }
}
</script>

<style scoped>
  .el-tag + .el-tag {
    margin-left: 10px;
  }
</style>
