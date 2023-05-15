<template>
  <Card @click="goToDialog(info.friend_id)" class="mx-4 my-2">
    <template #title>
      <div class="dialog-header">
        <Avatar :image="info.avatar" shape="circle" size="45" />
        <h4 class="ml-2">{{ info.friend_name }}</h4>
      </div>
    </template>
    <template #content>
      <div class="dialog-message mx-6 my-1">
        <p>{{ info.last_message }}</p>
      </div>
    </template>
  </Card>
</template>

<script>
import Card from 'primevue/card';
import Avatar from 'primevue/avatar';

export default {
  components: {
    Avatar,
    Card
  },
  methods: {
    goToDialog(id) {
      this.$router.push({ name: 'DialogPage', params: { id } });
    },
  },
  props: {
    info: {
      required: true,
      type: Object
    }
  }
}
</script>

<style scoped>
.dialog-header {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.dialog-message {
  text-align: left;
}
</style>