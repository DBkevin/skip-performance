import request from '../utils/request'

export const getProjectList = (params) => {
  return request({
    url: '/projects',
    method: 'get',
    params
  })
}

export const getProject = (id) => {
  return request({
    url: `/projects/${id}`,
    method: 'get'
  })
}

export const createProject = (data) => {
  return request({
    url: '/projects',
    method: 'post',
    data
  })
}

export const updateProject = (id, data) => {
  return request({
    url: `/projects/${id}`,
    method: 'put',
    data
  })
}

export const deleteProject = (id) => {
  return request({
    url: `/projects/${id}`,
    method: 'delete'
  })
}
