<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-10">
        <DataTable :value="people" responsive-layout="autoscroll" row-hover="true" class="p-datatable fixed-height hide-column-header">
          <Column>
            <template #body="slotProps">
              {{slotProps.data.name + " " + slotProps.data.surname}}
            </template>
          </Column>
          <Column field="interests"></Column>
          <Column>
            <template #body="slotProps">
              <Button v-if="!slotProps.data.is_friend && !slotProps.data.is_request_sent"
                class="p-button-sm" @click="onFriendshipRequest(slotProps.data, $event)">
                Friendship request
              </Button>
              <span v-if="slotProps.data.is_request_sent">Friendship request sent</span>
              <span v-if="slotProps.data.is_friend">Your friend</span>
            </template>
          </Column>
        </DataTable>
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
        axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/getpeople`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        }).then((res) => {
          // this.people = res.data;
          res.data.forEach((item) => {
              axios.get(`${process.env.VUE_APP_API_HOST || ''}/api/account/profile/${item.id}`, {
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
        axios.post(`${process.env.VUE_APP_API_HOST || ''}/api/account/friend_request`,
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
.hide-column-header table thead tr[role='row'] {
  display: none;
}
.fixed-height table tbody tr{
    height : 70px;
}
</style>
