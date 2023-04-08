<template>
  <h5>My info</h5>
  <table class="table">
    <tbody>
      <tr>
        <td>Name</td>
        <td>{{ profile.name }} {{ profile.surname }}</td>
      </tr>
      <tr>
        <td>Age</td>
        <td>{{ profile.age }}</td>
      </tr>
      <tr>
        <td>Sex</td>
        <td>{{ profile.sex }}</td>
      </tr>
      <tr>
        <td>City</td>
        <td>{{ profile.city }}</td>
      </tr>
      <tr>
        <td>Interests</td>
        <td>{{ profile.interests }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
export default {
  props: {
    profile: {
      required: true,
      type: Object,
      default: NaN,
    },
  },
};
</script>
