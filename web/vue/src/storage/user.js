class User {
  get () {
    return {
      'token': this.getToken(),
      'uid': this.getUid(),
      'username': this.getUsername(),
      'role': this.getRole(),
      'isAdmin': this.getIsAdmin(),
      'isDeveloper': this.getIsDeveloper()
    }
  }

  getToken () {
    return localStorage.getItem('token') || ''
  }

  setToken (token) {
    localStorage.setItem('token', token)
    return this
  }

  clear () {
    localStorage.clear()
  }

  getUid () {
    return localStorage.getItem('uid') || ''
  }

  setUid (uid) {
    localStorage.setItem('uid', uid)
    return this
  }

  getUsername () {
    return localStorage.getItem('username') || ''
  }

  setUsername (username) {
    localStorage.setItem('username', username)
    return this
  }

  setRole (role) {
    localStorage.setItem('role', role)
    return this
  }

  getRole () {
    let role = localStorage.getItem('role')
    if (role) {
      return parseInt(role)
    }
    return -1
  }

  getIsAdmin () {
    return this.getRole() === 1
  }

  getIsDeveloper () {
    return this.getRole() > 0
  }
}

export default new User()
