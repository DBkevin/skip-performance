import request from '../utils/request'

export const getVisitList = (params) => {
  return request({
    url: '/visits',
    method: 'get',
    params
  })
}

export const getVisit = (id) => {
  return request({
    url: `/visits/${id}`,
    method: 'get'
  })
}

export const createVisit = (data) => {
  return request({
    url: '/visits',
    method: