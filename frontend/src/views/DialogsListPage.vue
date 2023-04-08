<template>
  <div class="dialogs-list">
    <h2>Список диалогов</h2>
    <Card v-for="(dialog, index) in dialogs" :key="index" class="mx-4 my-2" @click="goToDialog(dialog.friend_id)">
      <template #title>
        <div class="dialog-header">
          <Avatar :image="dialog.avatar" shape="circle" size="45" />
          <h4 class="ml-2">{{ friends.get(dialog.friend_id) }}</h4>
        </div>
      </template>
      <template #content>
        <div class="dialog-message mx-6 my-1">
          <p>{{ dialog.last_message }}</p>
        </div>
      </template>
    </Card>
  </div>
</template>

<script>
import Card from 'primevue/card';
import Avatar from 'primevue/avatar';

import axios from 'axios';

export default {
  name: 'DialogsList',
  components: {
    Card,
    Avatar,
  },
  data() {
    return {
      dialogs: [],
      friends: new Map(),
    };
  },
  methods: {
    goToDialog(id) {
      this.$router.push({ name: 'DialogPage', params: { id } });
      // this.$router.push('/dialog/' + id)
    },
    getDialogList() {
      axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/dialog/list`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        }
      }).then((res) => {
        res.data.forEach((item) => {
          if (!this.friends.get(item.friend_id)) {
            axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/profile/${item.friend_id}`, {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
              },
            })
              .then((res) => {
                this.friends.set(res.data.id, `${res.data.name} ${res.data.surname}`)
              })
          }
        })

        this.dialogs = res.data
      })
    },
  },
  created() {
    this.getDialogList();
  },
};
</script>

<style scoped>
.dialog-header {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.dialogs-list h2 {
  margin-bottom: 1rem;
}

.dialog-message {
  text-align: left;
}
</style>