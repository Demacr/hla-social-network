<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-10">
        <table class="table">
          <tbody>
            <tr v-for="(profile, index) in people" :key="index">
              <td>{{profile.surname}} {{profile.name}}</td>
              <td>{{profile.age}}</td>
              <td>{{profile.city}}</td>
              <td>
                <button v-if="!profile.is_friend && !profile.is_request_sent"
                  type="button" class="btn btn-sm" @click="onFriendshipRequest(profile, $event)">
                  Friendship request
                </button>
                <span v-if="profile.is_request_sent">Friendship request sent</span>
                <span v-if="profile.is_friend">Your friend</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from 'axios';

  export default{
    data() {
      return {
        people: [],
      };
    },
    methods: {
      getPeople() {
        axios.get(`${process.env.API_HOST || ''}/api/account/getpeople`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        }).then((res) => {
          // this.people = res.data;
          res.data.forEach((item) => {
              axios.get(`${process.env.API_HOST || ''}/api/account/profile/${item.id}`, {
                headers: {
                  Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
              }).then((prof) => {
                const temp = item;
                temp.is_friend = prof.data.is_friend;
                temp.is_request_sent = prof.data.is_request_sent;
                this.people.push(item);
              })
              .catch((error) => {
                // eslint-disable-next-line
                console.log(error);
              });
            },
          );
        })
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        });
      },
      sendFriendshipRequest(payload) {
        axios.post(`${process.env.API_HOST || ''}/api/account/friend_request`,
          {
            friendID: payload.id,
          },
          {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        }).then(() => {
          payload.is_request_sent = true;
        })
        .catch((error) => {
          // eslint-disable-next-line
          console.log(error);
        });
      },
      onFriendshipRequest(profile, evt) {
        evt.preventDefault();
        this.sendFriendshipRequest(profile, evt);
      },
    },
    created() {
      this.getPeople();
    },
  };
</script>

<style>

</style>
