import request from '../utils/request'

export const getEmployeeList = (params) => {
  return request({
    url: '/employees',
    method: 'get',
    params
  })
}

export const getEmployee = (id) => {
  return request({
    url: `/employees/${id}`,
    method: 'get'
  })
}

export const createEmployee = (data) => {
  return request({
    url: '/employees',
    method: 'post',
    data
  })
}

export const updateEmployee = (id, data) => {
  return request({
    url: `/employees/${id}`,
    method: 'put',
    data
  })
}

export const deleteEmployee = (id) => {
  return request({
    url: `/employees/${id}`,
    method: 'delete'
  })
}
