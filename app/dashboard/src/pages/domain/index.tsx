import {
  Button,
  Card,
  Col,
  Divider,
  Form,
  Input,
  Row,
  message,
  Popconfirm,
  InputNumber,
} from 'antd';
import React, { Component, Fragment } from 'react';

import { Dispatch, Action } from 'redux';
import { FormComponentProps } from 'antd/es/form';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import { SorterResult } from 'antd/es/table';
import { connect } from 'dva';
import { StateType } from './model';
import CreateForm from './components/CreateForm';
import StandardTable, { StandardTableColumnProps } from './components/StandardTable';
// eslint-disable-next-line
import { TableListItem, TableListPagination, TableListParams } from './data';

import styles from './style.less';

const FormItem = Form.Item;
const getValue = (obj: { [x: string]: string[] }) =>
  Object.keys(obj)
    .map(key => obj[key])
    .join(',');

interface TableListProps extends FormComponentProps {
  dispatch: Dispatch<Action>;
  loading: boolean;
  domain: StateType;
}

interface TableListState {
  modalVisible: boolean;
  isEdit: boolean;
  currentRow: any;
  selectedRows: TableListItem[];
  formValues: { [key: string]: string };
}

/* eslint react/no-multi-comp:0 */
@connect(
  ({
    domain,
    loading,
  }: {
    domain: StateType;
    loading: {
      models: {
        [key: string]: boolean;
      };
    };
  }) => ({
    domain,
    loading: loading.models.domain,
  }),
)
class TableList extends Component<TableListProps, TableListState> {
  state: TableListState = {
    modalVisible: false,
    isEdit: false,
    currentRow: {},
    selectedRows: [],
    formValues: {},
  };

  columns: StandardTableColumnProps[] = [
    {
      title: '域名ID',
      dataIndex: 'id',
    },
    {
      title: '域名',
      dataIndex: 'domain',
    },
    {
      title: '米表ID',
      dataIndex: 'panel_id',
    },
    {
      title: '分类ID',
      dataIndex: 'cat_id',
    },
    {
      title: '简介',
      dataIndex: 'desc',
    },
    {
      title: '购入成本',
      dataIndex: 'cost',
    },
    {
      title: '续费成本',
      dataIndex: 'renew',
    },
    {
      title: '购入时间',
      dataIndex: 'buy',
    },
    {
      title: '注册商',
      dataIndex: 'registrar',
    },
    {
      title: '注册时间',
      dataIndex: 'create',
    },
    {
      title: '到期时间',
      dataIndex: 'expire',
    },
    {
      title: '管理操作',
      render: (text, record) => (
        <Fragment>
          <a
            onClick={() =>
              this.setState(prevState => ({
                ...prevState,
                currentRow: record,
                isEdit: true,
                modalVisible: true,
              }))
            }
          >
            修改
          </a>
          <Divider type="vertical" />
          <Popconfirm
            title={`确认删除域名「${record.domain}」`}
            onConfirm={() => {
              this.handleDelete(record);
            }}
            okText="Yes"
            cancelText="No"
          >
            <a>删除</a>
          </Popconfirm>
        </Fragment>
      ),
    },
  ];

  componentDidMount() {
    const { dispatch } = this.props;
    const { formValues } = this.state;

    dispatch({
      type: 'domain/fetch',
      payload: formValues,
    });
  }

  handleDelete = (record: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'domain/remove',
      payload: record,
      callback: () => {
        dispatch({
          type: 'domain/fetch',
          payload: this.state.formValues,
        });
        message.success('删除成功');
      },
    });
  };

  handleStandardTableChange = (
    pagination: Partial<TableListPagination>,
    filtersArg: Record<keyof TableListItem, string[]>,
    sorter: SorterResult<TableListItem>,
  ) => {
    const { dispatch } = this.props;
    const { formValues } = this.state;

    const filters = Object.keys(filtersArg).reduce((obj, key) => {
      const newObj = { ...obj };
      newObj[key] = getValue(filtersArg[key]);
      return newObj;
    }, {});

    const params: Partial<TableListParams> = {
      currentPage: pagination.current,
      pageSize: pagination.pageSize,
      ...formValues,
      ...filters,
    };
    if (sorter.field) {
      params.sorter = `${sorter.field}_${sorter.order}`;
    }

    dispatch({
      type: 'domain/fetch',
      payload: params,
    });
  };

  handleFormReset = () => {
    const { form, dispatch } = this.props;
    form.resetFields();
    this.setState({
      formValues: {},
    });
    dispatch({
      type: 'domain/fetch',
      payload: {},
    });
  };

  handleSelectRows = (rows: TableListItem[]) => {
    this.setState({
      selectedRows: rows,
    });
  };

  handleSearch = (e: React.FormEvent) => {
    e.preventDefault();

    const { dispatch, form } = this.props;

    form.validateFields((err, fieldsValue) => {
      if (err) return;

      const values = {
        ...fieldsValue,
        updatedAt: fieldsValue.updatedAt && fieldsValue.updatedAt.valueOf(),
      };

      this.setState({
        formValues: values,
      });

      dispatch({
        type: 'domain/fetch',
        payload: values,
      });
    });
  };

  handleModalVisible = (flag?: boolean) => {
    this.setState({
      modalVisible: !!flag,
      currentRow: {},
      isEdit: false,
    });
  };

  handleAdd = (fields: any, isEdit: boolean) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'domain/add',
      payload: fields,
      callback: () => {
        dispatch({
          type: 'domain/fetch',
          payload: this.state.formValues,
        });
        message.success(`${isEdit ? '修改' : '添加'}成功`);
        this.handleModalVisible();
      },
    });
  };

  renderSimpleForm() {
    const { form } = this.props;
    const { getFieldDecorator } = form;
    return (
      <Form onSubmit={this.handleSearch} layout="inline">
        <Row gutter={{ md: 8, lg: 24, xl: 48 }}>
          <Col md={8} sm={24}>
            <FormItem label="域名称">
              {getFieldDecorator('domain')(<Input placeholder="请输入" />)}
            </FormItem>
          </Col>
          <Col md={4} sm={12}>
            <FormItem label="米表ID">
              {getFieldDecorator('panel_id')(<InputNumber min={1} />)}
            </FormItem>
          </Col>
          <Col md={4} sm={12}>
            <FormItem label="分类ID">
              {getFieldDecorator('cat_id')(<InputNumber min={1} />)}
            </FormItem>
          </Col>
          <Col md={8} sm={24}>
            <span className={styles.submitButtons}>
              <Button type="primary" htmlType="submit">
                查询
              </Button>
              <Button style={{ marginLeft: 8 }} onClick={this.handleFormReset}>
                重置
              </Button>
            </span>
          </Col>
        </Row>
      </Form>
    );
  }

  render() {
    const {
      domain: { data },
      loading,
      dispatch,
    } = this.props;

    const { selectedRows, modalVisible, isEdit, currentRow } = this.state;

    const parentMethods = {
      handleAdd: this.handleAdd,
      handleModalVisible: this.handleModalVisible,
    };

    return (
      <PageHeaderWrapper>
        <Card bordered={false}>
          <div className={styles.tableList}>
            <div className={styles.tableListForm}>{this.renderSimpleForm()}</div>
            <div className={styles.tableListOperator}>
              <Button icon="plus" type="primary" onClick={() => this.handleModalVisible(true)}>
                新建
              </Button>
            </div>
            <StandardTable
              rowKey="id"
              scroll={{ x: 1500 }}
              selectedRows={selectedRows}
              loading={loading}
              data={data}
              columns={this.columns}
              onSelectRow={this.handleSelectRows}
              onChange={this.handleStandardTableChange}
            />
          </div>
        </Card>
        <CreateForm
          {...parentMethods}
          dispatch={dispatch}
          currentRow={currentRow}
          isEdit={isEdit}
          modalVisible={modalVisible}
        />
      </PageHeaderWrapper>
    );
  }
}

export default Form.create<TableListProps>()(TableList);
