new Vue({
  el: '#app',
  data() {
    return {
      plugins: null,
      errored: false
    }
  },
  mounted() {
    axios
      .get('http://localhost:8080/feed/list')
      .then(
        response => {
          this.plugins = response.data
        }
      )
      .catch(error => {
        console.log(error)
        this.errored = true
      })
  }
})