import request from '../utils/request'

export const getPerformanceReport = (params) => {
  return request({
    url: '/reports/performance',
    method: 'get',
    params
  })
}

export const getEmployeePerformance = (params) => {
  return request({
    url: '/reports/employee-performance',
    method: 'get',
    params
  })
}

export const getProjectPerformance = (params) => {
  return request({
    url: '/reports/project-performance',
    method: 'get',
    params
  })
}
