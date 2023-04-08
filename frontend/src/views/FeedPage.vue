<template>
  <div class="feed">
    <h2>Лента постов</h2>
    <Card v-for="(post, index) in posts" :key="index" class="mx-4 my-4">
      <template #title>
        <div class="post-header">
          <!-- <Avatar :image="post.profile_id" shape="circle" size="45" /> -->
          <h4 class="ml-2">{{ authors.get(post.profile_id) }}</h4>
        </div>
      </template>
      <template #content>
        <h4>{{ post.title }}</h4>
        <p>{{ post.text }}</p>
      </template>
    </Card>
  </div>
</template>

<script>
import Card from 'primevue/card';
// import Avatar from 'primevue/avatar';

import axios from 'axios';

export default {
  name: 'FeedPage',
  components: {
    Card,
    // Avatar,
  },
  data() {
    return {
      postIds: [],
      posts: [],
      authors: new Map(),
    };
  },
  methods: {
    getFeed() {
      axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/post/feed`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      }).then((res) => {
        this.postIds = res.data.slice(0, 20)
        this.posts = []

        this.postIds.forEach((item) => {
          axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/post/${item}`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            }
          })
            .then((res) => {
              if (!this.authors.get(res.data.profile_id)) {
                axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/profile/${res.data.profile_id}`, {
                  headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                  },
                })
                  .then((res) => {
                    this.authors.set(res.data.id, `${res.data.name} ${res.data.surname}`)
                  })
              }

              this.posts.push(res.data)
            })
            .catch((error) => {
              // eslint-disable-next-line
              console.log(error);
            })
        });
      })
    },
  },
  created() {
    this.getFeed();
  },
};
</script>

<style scoped>
.post-header {
  display: flex;
  align-items: center;
}

.feed h2 {
  margin-bottom: 1rem;
}
</style>
