<template>
  <div class="dialog sm:w-30rem md:w-9 xl:w-7">
    <h2>Диалог с {{ interlocutor.name }} {{ interlocutor.surname }}</h2>
    <div class="messages">
      <div v-for="(message, index) in messages" :key="index"
        :class="['message', 'mx-4', message.From == this.$props.id ? 'received' : 'sent']">
        <div class="message-content">{{ message.text }}</div>
        <div class="message-timestamp">{{ formatTimestamp(message.Timestamp) }}</div>
      </div>
      <div ref="last-message">
      </div>
    </div>
    <div class="message-input mx-4 mt-auto mb-4">
      <InputText v-model="newMessage" placeholder="Введите сообщение..." />
      <Button label="Отправить" @click="sendMessage" />
    </div>
  </div>
</template>

<script>
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';

import axios from 'axios';

export default {
  name: 'DialogPage',
  components: {
    InputText,
    Button,
  },
  data() {
    return {
      interlocutor: {
        name: 'Иван Иванов',
        avatar: 'https://example.com/avatar1.jpg',
      },
      messages: [
        {
          sender: 'me',
          content: 'Привет, как дела?',
        },
        {
          sender: 'interlocutor',
          content: 'Привет! Всё хорошо, спасибо!',
        },
        // Другие сообщения
      ],
      newMessage: '',
    };
  },
  methods: {
    sendMessage() {
      if (this.newMessage.trim() === '') return;

      axios.post(`${process.env.VUE_APP_API_HOST || ''}/api/dialog/${this.$props.id}/send`, {
        text: this.newMessage.trim()
      },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          }
        }).then(() => {
          this.getDialog();
        })

      this.newMessage = '';
    },
    getDialog() {
      axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/dialog/${this.$props.id}/list`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        }
      }).then((res) => {
        this.messages = res.data;
      });
    },
    getInterlocutorName() {
      axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/profile/${this.$props.id}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      }).then((res) => {
        this.interlocutor = res.data;
      })
    },
    formatTimestamp(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    },

    goToLastMessage() {
      this.$nextTick().then(() =>
        this.$refs['last-message'].scrollIntoView({ behavior: "smooth" })
      )
    },
  },

  created() {
    this.getInterlocutorName();
    this.getDialog();
    this.goToLastMessage();
  },

  watch: {
    messages(nm) {
      if (nm) {
        this.goToLastMessage();
      }
    }
  },

  props: {
    id: {
      type: Number,
      default: 0
    },
  },
};
</script>

<style scoped>
.dialog {
  display: flex;
  flex-direction: column;
  margin: auto;
  height: calc(100vh - 108px);
}

.messages {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1rem;
  overflow-y: auto;
}

.message {
  max-width: 50%;
  padding: 1rem;
  border-radius: 8px;
}

.sent {
  align-self: flex-end;
  background-color: #f5f5f5;
  text-align: right;
}

.received {
  align-self: flex-start;
  background-color: #e5e5e5;
  text-align: left;
}

.message-input {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-end;
}

.message-timestamp {
  font-size: 0.5rem;
  color: gray;
}
</style>
