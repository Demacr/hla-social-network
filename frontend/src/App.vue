<template>
  <nav>
    <router-link to="/">Homepage</router-link> |
    <router-link to="/feed">Feed</router-link> |
    <router-link to="/dialogs">Dialogs</router-link> |
    <router-link to="/people">People</router-link>
    <span v-if="isLoggedIn"> |
      <router-link to="/logout">Logout</router-link>
    </span>
  </nav>
  <br>
  <router-view />
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
        const path = `${process.env.VUE_APP_API_HOST || ''}/api/account/myinfo`;
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
            // if (error.response.status === 401) {
            this.$router.push('/login');
            // }
          },
          );
      } else {
        this.$router.push('/login');
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
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

nav {
  padding: 30px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
}

nav a.router-link-exact-active {
  color: #42b983;
}
</style>
