<template>
    <div class='container'>
        <form>
            <div class="form-group">
                <label for="exampleInputEmail1">Email address</label>
                <input type="email" class="form-control" id="exampleInputEmail1"
                    aria-describedby="emailHelp" placeholder="Enter email"
                    v-model="authorizeForm.email"
                    required>
                <small id="emailHelp" class="form-text text-muted">
                    We'll never share your email with anyone else.
                </small>
            </div>
            <div class="form-group">
                <label for="exampleInputPassword1">Password</label>
                <input type="password" class="form-control" id="exampleInputPassword1"
                    placeholder="Password"
                    v-model="authorizeForm.password"
                    required>
            </div>
            <button type="submit" class="btn btn-primary btn-sm" @click="onAuthorize">
                Sign In
            </button>
            <button type="button" class="btn btn-success btn-sm" v-b-modal.registrate-modal>
                Sign Up
            </button>
        </form>

        <b-modal ref="registrateModal"
                id="registrate-modal"
                title="Registrate a new user"
                hide-footer>
            <b-form @submit="onSubmit" @reset="onReset" class="w-100">
                <b-form-group id="form-name-group"
                                label="Name:"
                                label-for="form-name-input">
                    <b-form-input id="form-name-input"
                                type="text"
                                v-model="registrateForm.name"
                                required
                                placeholder="Enter name">
                    </b-form-input>
                </b-form-group>
                <b-form-group id="form-surname-group"
                                label="Surname:"
                                label-for="form-surname-input">
                    <b-form-input id="form-surname-input"
                                type="text"
                                v-model="registrateForm.surname"
                                required
                                placeholder="Enter surname">
                    </b-form-input>
                </b-form-group>
                <b-form-group id="form-age-group"
                                label="Age:"
                                label-for="form-age-input">
                    <b-form-input id="form-age-input"
                                type="number"
                                v-model.number="registrateForm.age"
                                required
                                placeholder="Enter age">
                    </b-form-input>
                </b-form-group>
                <b-form-group id="form-sex-group"
                                label="Sex:"
                                label-for="form-sex-input">
                    <b-form-radio-group>
                        <b-form-radio
                                    type="radio"
                                    name="sex-radio"
                                    v-model="registrateForm.sex"
                                    id="form-sex-input1"
                                    required
                                    value="Male">
                            Male
                        </b-form-radio>
                        <b-form-radio
                                    type="radio"
                                    name="sex-radio"
                                    v-model="registrateForm.sex"
                                    id="form-sex-input2"
                                    required
                                    value="Female">
                            Female
                        </b-form-radio>
                    </b-form-radio-group>
                </b-form-group>
                <b-form-group id="form-interests-group"
                                label="Interests:"
                                label-for="form-interests-input">
                    <b-form-input id="form-interests-input"
                                type="text"
                                v-model="registrateForm.interests"
                                required
                                placeholder="Enter interests, comma-separated">
                    </b-form-input>
                </b-form-group>
                <b-form-group id="form-city-group"
                                label="City:"
                                label-for="form-city-input">
                    <b-form-input id="form-city-input"
                                type="text"
                                v-model="registrateForm.city"
                                required
                                placeholder="Enter city">
                    </b-form-input>
                </b-form-group>
                <b-form-group id="form-email-group"
                                label="Email:"
                                label-for="form-email-input">
                    <b-form-input id="form-email-input"
                                type="email"
                                v-model="registrateForm.email"
                                required
                                placeholder="Enter email">
                    </b-form-input>
                </b-form-group>
                <b-form-group id="form-password-group"
                                label="Password:"
                                label-for="form-password-input">
                    <b-form-input id="form-name-input"
                                type="password"
                                v-model="registrateForm.password"
                                required
                                placeholder="Enter password">
                    </b-form-input>
                </b-form-group>
                <b-button type="submit" variant="primary">Submit</b-button>
                <b-button type="reset" variant="danger">Reset</b-button>
            </b-form>
        </b-modal>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'Index',
    data() {
        return {
            registrateForm: {
                name: null,
                surname: null,
                age: null,
                sex: null,
                interests: null,
                city: null,
                email: null,
                password: null,
            },
            authorizeForm: {
                email: null,
                password: null,
            },
        };
    },
    methods: {
        authorize(payload) {
            const path = 'http://localhost:1323/api/authorize';
            axios.post(path, payload)
                .then(() => {
                    this.$router.push('/home');
                })
                .catch((error) => {
                    // eslint-disable-next-line
                    console.log(error);
                },
                );
        },
        registrate(payload) {
            const path = 'http://localhost:1323/api/registrate';
            axios.post(path, payload)
                .catch((error) => {
                    // eslint-disable-next-line
                        console.log(error);
                    },
                );
        },
        initForm() {
            this.registrateForm.name = null;
            this.registrateForm.surname = null;
            this.registrateForm.age = null;
            this.registrateForm.sex = null;
            this.registrateForm.interests = null;
            this.registrateForm.city = null;
            this.registrateForm.email = null;
            this.registrateForm.password = null;
        },
        onAuthorize(evt) {
            evt.preventDefault();
            const payload = this.authorizeForm;
            this.authorize(payload);
        },
        onSubmit(evt) {
            evt.preventDefault();
            const payload = this.registrateForm;
            this.registrate(payload);
            this.$refs.registrateModal.hide();
        },
        onReset(evt) {
            evt.preventDefault();
            this.$refs.registrateModal.hide();
            this.initForm();
        },
    },
};
</script>

<style>

</style>
