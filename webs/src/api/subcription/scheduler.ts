import request from '@/utils/request';

// 订阅调度器接口定义
export interface SubScheduler {
  ID: number;
  Name: string;
  URL: string;
  CronExpr: string;
  Enabled: boolean;
}

export interface SubSchedulerRequest {
  id?: number;
  name: string;
  url: string;
  cron_expr: string;
  enabled: boolean;
}

// 获取订阅调度器列表
export function getSubSchedulers() {
  return request({
    url: '/api/v1/sub_scheduler/get',
    method: 'get'
  });
}

// 添加订阅调度器
export function addSubScheduler(data: SubSchedulerRequest) {
  return request({
    url: '/api/v1/sub_scheduler/add',
    method: 'post',
    data
  });
}

// 更新订阅调度器
export function updateSubScheduler(data: SubSchedulerRequest) {
  return request({
    url: '/api/v1/sub_scheduler/update',
    method: 'put',
    data
  });
}

// 删除订阅调度器
export function deleteSubScheduler(id: number) {
  return request({
    url: `/api/v1/sub_scheduler/delete/${id}`,
    method: 'delete'
  });
}
