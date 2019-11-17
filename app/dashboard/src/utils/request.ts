import axios, { AxiosError } from 'axios';
import { notification } from 'antd';

const inst = axios.create({
  timeout: 20000,
  withCredentials: true,
  headers: {},
});

// @cc: 检测 axios 响应状态
function onStatusError(error: AxiosError | Error) {
  const err =
    'response' in error && error.response
      ? {
          code: error.response.status,
          message: error.response.statusText,
        }
      : { code: 9999, message: error.message };
  if (err.code === 401 || err.code === 403) {
    // @todo 未登录未授权
    // EventCenter.emit('common.user.status', err);
    notification.error({
      message: '登录信息失效',
      description: '即将跳转到登录界面',
    });
    localStorage.removeItem('nbdomain-token');
    localStorage.removeItem('nbdomain-token-expired');
    window.location.href = '/user/login';
  }
  notification.error({
    message: `请求错误 ${err.code}`,
    description: err.message,
  });
  return false;
}

export type AjaxPromise<R> = Promise<R>;

export interface ExtraFetchParams {
  extra?: any;
}

export interface WrappedFetchParams extends ExtraFetchParams {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'OPTIONS' | 'PATCH' | 'HEAD';
  url: string;
  data?: any; // post json
  form?: any; // post form
  query?: any;
  header?: any;
  path?: any;
}

export class WrappedFetch {
  /**
   * @description ajax 方法
   */
  // eslint-disable-next-line
  public async ajax(
    { method, url, data, form, query, header, extra }: WrappedFetchParams,
    path?: string,
    basePath?: string,
  ) {
    let config = {
      ...extra,
      method: method.toLocaleLowerCase(),
      headers: { ...header },
    };

    // 授权 token
    const token = localStorage.getItem('nbdomain-token');
    const expired = localStorage.getItem('nbdomain-token-expired');
    if (token && expired && Date.parse(expired) > new Date().getMilliseconds()) {
      config.headers.Authorization = `token ${token}`;
    }

    // json
    if (data) {
      config = {
        ...config,
        headers: {
          ...config.headers,
          'Content-Type': 'application/json',
        },
        data,
      };
    }
    // form
    if (form) {
      const postData = new FormData();
      Object.keys(form).forEach(k => {
        if (form[k]) postData.append(k, form[k]);
      });
      config = {
        ...config,
        headers: {
          ...config.headers,
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        data: postData,
      };
    }
    return inst
      .request({ ...config, url, params: query })
      .then((res: { data: any }) => {
        if (res.data.code !== 200) {
          const err = {
            response: {
              status: res.data.code,
              statusText: res.data.message,
            },
          };
          console.debug('ajax', path, basePath);
          throw err;
        } else if (res.data.message) {
          notification.success({
            message: `请求成功 ${res.data.code}`,
            description: res.data.message,
          });
        }
        return res.data.result;
      })
      .catch(onStatusError);
  }

  /**
   * @description 接口传参校验
   */
  // eslint-disable-next-line
  public check<V>(value: V, name: string) {
    if (value === null || value === undefined) {
      const msg = `[ERROR PARAMS]: ${name} can't be null or undefined`;
      // 非生产环境，直接抛出错误
      if (process.env.NODE_ENV === 'development') {
        throw Error(msg);
      }
    }
  }
}

export default new WrappedFetch();
