<template>
  <div class="dialogs-list">
    <h2>Список диалогов</h2>
    <dialog-item v-for="(dialog, index) in dialogs" :key="dialog.dialog_id" :info="dialog"></dialog-item>
  </div>
</template>

<script>
import axios from 'axios';

import DialogItem from '@/components/DialogItem.vue'

export default {
  name: 'DialogsList',
  components: {
    DialogItem,
  },
  data() {
    return {
      dialogs: [],
    };
  },
  methods: {
    getDialogList() {
      axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/dialog/list`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        }
      }).then((res) => {
        this.dialogs = res.data;
        this.dialogs.forEach((item) => {
          axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/profile/${item.friend_id}`, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
          })
            .then((res) => {
              item.friend_name = `${res.data.name} ${res.data.surname}`;
            })
        })
      })
    },
  },
  created() {
    this.getDialogList();
  },
};
</script>

<style scoped>
.dialogs-list h2 {
  margin-bottom: 1rem;
}
</style>