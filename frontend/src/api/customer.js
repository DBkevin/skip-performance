import request from '../utils/request'

export const getCustomerList = (params) => {
  return request({
    url: '/customers',
    method: 'get',
    params
  })
}

export const getCustomer = (id) => {
  return request({
    url: `/customers/${id}`,
    method: 'get'
  })
}

export const createCustomer = (data) => {
  return request({
    url: '/customers',
    method: 'post',
    data
  })
}

export const updateCustomer = (id, data) => {
  return request({
    url: `/customers/${id}`,
    method: 'put',
    data
  })
}

export const deleteCustomer = (id) => {
  return request({
    url: `/customers/${id}`,
    method: 'delete'
  })
}
