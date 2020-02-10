import { shallowMount } from '@vue/test-utils'
import UrlUploadModal from '@/components/modals/url-upload-modal'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(UrlUploadModal, {
    stubs: ['b-modal', 'b-form', 'b-form-input', 'b-form-text', 'b-button'],
    propsData: {
      id: 'id-abc-123',
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('UrlUploadModal', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(UrlUploadModal)).toBe(true)
  })
})
