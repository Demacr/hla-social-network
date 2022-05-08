<template>
  <div id="app">
    <div id="nav">
      <router-link to="/home">Homepage</router-link>
      <router-link to="/people">People</router-link>
      <span v-if="isLoggedIn" class="float-right">
        <router-link to="/logout">Logout</router-link>
      </span>
    </div>
    <router-view/>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      isLoggedIn: false,
    };
  },
  methods: {
    checkToken() {
      if (localStorage.getItem('token')) {
        const path = `${process.env.API_HOST || ''}/api/account/myinfo`;
        axios.get(
          path,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
          },
        )
          .then(() => {
            this.isLoggedIn = true;
          })
          .catch((error) => {
            // eslint-disable-next-line
            console.log(error);
            if (error.response.status === 401) {
              this.$router.push('/');
            }
          },
          );
      } else {
        this.$router.push('/');
      }
    },
    logout() {

    },
  },
  created() {
    this.checkToken();
  },
};
</script>

<style>
#app {
  margin-top: 60px;
}
</style>
