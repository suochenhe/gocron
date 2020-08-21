import httpClient from '../utils/httpClient'

export default {
  ding (callback) {
    httpClient.get('/system/ding', {}, callback)
  },
  updateDing (data, callback) {
    httpClient.post('/system/ding/update', data, callback)
  },
  createDingUser (name, mobile, callback) {
    httpClient.post('/system/ding/user', {name, mobile}, callback)
  },
  removeDingUser (userId, callback) {
    httpClient.post(`/system/ding/user/remove/${userId}`, {}, callback)
  },
  slack (callback) {
    httpClient.get('/system/slack', {}, callback)
  },
  updateSlack (data, callback) {
    httpClient.post('/system/slack/update', data, callback)
  },
  createSlackChannel (channel, callback) {
    httpClient.post('/system/slack/channel', {channel}, callback)
  },
  removeSlackChannel (channelId, callback) {
    httpClient.post(`/system/slack/channel/remove/${channelId}`, {}, callback)
  },
  mail (callback) {
    httpClient.get('/system/mail', {}, callback)
  },
  updateMail (data, callback) {
    httpClient.post('/system/mail/update', data, callback)
  },
  createMailUser (data, callback) {
    httpClient.post('/system/mail/user', data, callback)
  },
  removeMailUser (userId, callback) {
    httpClient.post(`/system/mail/user/remove/${userId}`, {}, callback)
  },
  webhook (callback) {
    httpClient.get('/system/webhook', {}, callback)
  },
  updateWebHook (data, callback) {
    httpClient.post('/system/webhook/update', data, callback)
  }
}
