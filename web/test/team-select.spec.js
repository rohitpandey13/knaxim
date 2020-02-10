import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'

import TeamSelect from '@/components/team-select'
import { ACTIVATE_GROUP } from '@/store/mutations.type'
import { REFRESH_GROUPS } from '@/store/actions.type'

const localVue = createLocalVue()
localVue.use(Vuex)

let store = new Vuex.Store({
  actions: {
    [REFRESH_GROUPS] () {
      return true
    }
  },
  mutations: {
    [ACTIVATE_GROUP] () {
      return true
    }
  },
  getters: {
    currentUser () {
      return {
        id: 'user_id_abc_123'
      }
    }
  }
})

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(TeamSelect, {
    stubs: ['b-form-select'],
    store,
    localVue,
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      options () {
        return [{
          value: 'id-abc-123',
          text: 'group-name'
        }]
      },
      ...options.computed
    }
  })
}

describe('TeamSelect', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(TeamSelect)).toBe(true)
  })
})
