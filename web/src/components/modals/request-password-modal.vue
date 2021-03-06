<!--
// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
-->
<template>
  <b-modal
  :id="id"
  ref="modal"
  :no-close-on-backdrop="loading"
  :no-close-on-esc="loading"
  @hidden="close"
  centered
  hide-footer
  hide-header
  content-class="modal-style">
    <b-img src="@/assets/CloudEdison.png" alt="Cloud Edison"/>
    <b-form @submit.prevent="send">
      <b-form-group>
        <b-form-input autofocus v-model="name" placeholder="Username"/>
      </b-form-group>
      <div v-if="loading">
        <b-spinner class="m-3"/>
      </div>
      <b-form-group v-else>
        <b-button type="submit" class="shadow-sm" :disabled="name.length < 6">Send Reset Email</b-button>
        <b-button @click="toLogin" class="shadow-sm">Back to Login</b-button>
      </b-form-group>
    </b-form>
  </b-modal>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import { SEND_RESET_REQUEST } from '@/store/actions.type'

export default {
  name: 'request-password-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      name: ''
    }
  },
  computed: {
    ...mapGetters({
      loading: 'authLoading'
    })
  },
  methods: {
    close () {
      this.$emit('close')
    },
    send () {
      if (this.name.length < 6) {
        return
      }
      this.sendRequest({ name: this.name })
        .then(() => {
          this.toLogin()
        })
    },
    toLogin () {
      this.close()
    },
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.name = ''
      this.$refs['modal'].hide()
    },
    ...mapActions({ sendRequest: SEND_RESET_REQUEST })
  }
}
</script>

<style scoped lang="scss">

img {
  width: 50%;
}

input {
  margin-top: 10px;
  margin-bottom: 10px;
  width: 80%;
  display: inline-block;
}

button {
  @extend %pill-buttons;
  width: flex;
  margin-right: 5px;
  margin-left: 5px;
}

::v-deep .modal-style {
  @extend %modal-corners;
  text-align: center;
}

</style>
