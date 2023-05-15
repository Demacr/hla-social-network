<template>
    <div>
        <form>
            <div class="login">
                <div>
                    <label class="col-fixed" style="width:100px" for="username">Email</label>
                    <InputText id="username" type="text" v-model="authorizeForm.email"></InputText>
                </div>
                <div>
                    <label class="col-fixed" style="width:100px" for="password">Password</label>
                    <InputText id="password" type="password" v-model="authorizeForm.password"></InputText>
                </div>
                <div class="login-btns">
                    <div class="m-1">
                        <Button type="button" class="p-button-success" @click="onAuthorize">Login</Button>
                    </div>
                    <div class="m-1">
                        <Button type="button" @click="onRegistrate">Registrate</Button>
                    </div>
                </div>
            </div>
        </form>
    </div>
    <registration-form v-model:visible="showRegistrate"></registration-form>
</template>

<script>
import RegistrationForm from '@/components/RegistrationForm.vue';

import axios from 'axios';

export default {
    name: 'IndexPage',
    components: {
        RegistrationForm
    },
    data() {
        return {
            authorizeForm: {
                email: null,
                password: null,
            },
            showRegistrate: false,
        };
    },
    methods: {
        authorize(payload) {
            const path = `${process.env.VUE_APP_API_HOST || ''}/api/authorize`;
            // eslint-disable-next-line
            console.log(path);
            axios.post(path, payload)
                .then((result) => {
                    localStorage.setItem('token', result.data.token);
                    this.$router.push('/');
                })
                .catch(() => {
                },
                );
        },
        onAuthorize() {
            const payload = this.authorizeForm;
            this.authorize(payload);
        },
        onRegistrate() {
            this.showRegistrate = true;
        },
    },
};
</script>

<style>
.login {
    display: flex;
    flex-direction: column;
    row-gap: 1rem;
    align-items: flex-start;
    justify-content: baseline;
    text-align: left;
}

.login-btns {
    display: flex;
    flex-direction: row;
}
</style>
