<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-10">
        <h1>My page</h1>
        <hr><br><br>

        <my-info :profile="this.me"></my-info>

        <span v-if="isFrienshipRequests">
          <hr><br>
          <table class="table">
            <tbody>
              <tr v-for="(info, index) in friendInfo" :key="index">
                <td>
                  {{ info.name }} {{ info.surname }}
                </td>
                <td>
                  <button v-if="info.accepted === undefined" type="button" class="btn btn-sm"
                    @click="onFriendshipAccept(info.id, index, $event)">
                    Friendship accept
                  </button>
                  <span v-if="info.accepted === true">
                    Request accepted
                  </span>
                </td>
                <td>
                  <button v-if="info.accepted === undefined" type="button" class="btn btn-sm"
                    @click="onFriendshipDecline(info.id, index, $event)">
                    Friendship decline
                  </button>
                  <span v-if="info.accepted === false">
                    Request declined
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'HomePage',
  data() {
    return {
      me: {},
      isFrienshipRequests: false,
      friendRequests: {},
      friendInfo: [],
    };
  },
  methods: {
    getMyInfo() {
      const path = `${process.env.VUE_APP_API_HOST || ''}/api/account/myinfo`;
      return axios.get(
        path,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        },
      )
        .then((result) => {
          this.me = result.data;
        })
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        },
        );
    },
    getAccountInfo(id) {
      const path = `${process.env.VUE_APP_API_HOST || ''}/api/account/profile/${id}`;
      return axios.get(
        path,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        },
      )
        .then(result => result.data)
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        });
    },
    getMyFriendRequests() {
      const path = `${process.env.VUE_APP_API_HOST || ''}/api/account/my_friend_requests`;
      axios.get(
        path,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        },
      )
        .then((result) => {
          if (result.status === 200) {
            this.isFrienshipRequests = true;
            this.friendRequests = result.data;
          }
          this.friendRequests.map((request, i) => {
            this.getAccountInfo(request.friendID).then((response) => {
              this.$set(this.friendInfo, i, response);
            });
          });
        })
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        });
    },
    friendshipAccept(id, index) {
      const path = `${process.env.VUE_APP_API_HOST || ''}/api/account/friendship_request_accept`;
      return axios.post(
        path,
        {
          friendID: id,
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        },
      )
        .then(() => {
          this.friendInfo[index].accepted = true;
        })
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        });
    },
    friendshipDecline(id, index) {
      const path = `${process.env.VUE_APP_API_HOST || ''}/api/account/friendship_request_decline`;
      return axios.post(
        path,
        {
          friendID: id,
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        },
      )
        .then(() => {
          this.friendInfo[index].accepted = false;
        })
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        });
    },
    onFriendshipAccept(id, index, evt) {
      evt.preventDefault();
      this.friendshipAccept(id, index);
    },
    onFriendshipDecline(id, index, evt) {
      evt.preventDefault();
      this.friendshipDecline(id, index);
    },
  },
  created() {
    this.getMyInfo();
    this.getMyFriendRequests();
  },
};
</script>

<style></style>
