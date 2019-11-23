import { Button, Card, Col, Divider, Form, Input, Row, message, Popconfirm } from 'antd';
import React, { Component, Fragment } from 'react';

import { Dispatch, Action } from 'redux';
import { FormComponentProps } from 'antd/es/form';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import { SorterResult } from 'antd/es/table';
import { connect } from 'dva';
import { StateType } from './model';
import CreateForm from './components/CreateForm';
import StandardTable, { StandardTableColumnProps } from './components/StandardTable';
import { TableListItem, TableListPagination, TableListParams } from './data.d';

import styles from './style.less';
import ImportForm from './components/ImportForm';

const FormItem = Form.Item;
const getValue = (obj: { [x: string]: string[] }) =>
  Object.keys(obj)
    .map(key => obj[key])
    .join(',');

interface TableListProps extends FormComponentProps {
  dispatch: Dispatch<Action>;
  loading: boolean;
  panel: StateType;
}

interface TableListState {
  createModelVisible: boolean;
  importModalVisible: boolean;
  isEdit: boolean;
  currentPanel: any;
  currentRow: any;
  selectedRows: TableListItem[];
  formValues: { [key: string]: string };
}

/* eslint react/no-multi-comp:0 */
@connect(
  ({
    panel,
    loading,
  }: {
    panel: StateType;
    loading: {
      models: {
        [key: string]: boolean;
      };
    };
  }) => ({
    panel,
    loading: loading.models.panel,
  }),
)
class TableList extends Component<TableListProps, TableListState> {
  state: TableListState = {
    createModelVisible: false,
    importModalVisible: false,
    isEdit: false,
    currentPanel: {},
    currentRow: {},
    selectedRows: [],
    formValues: {},
  };

  columns: StandardTableColumnProps[] = [
    {
      title: '米表ID',
      dataIndex: 'id',
    },
    {
      title: '域名',
      dataIndex: 'domain',
    },
    {
      title: '标题[中]',
      dataIndex: 'name',
    },
    {
      title: '建仓成本',
      dataIndex: 'total_buy',
      render: (value, record) => <p>{record.total_buy ? value : '0'} 元</p>,
    },
    {
      title: '持仓成本',
      dataIndex: 'total_renew',
      render: (value, record) => <p>{record.total_renew ? value : '0'} 元</p>,
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
                createModelVisible: true,
              }))
            }
          >
            修改
          </a>
          <Divider type="vertical" />
          <a
            onClick={() => {
              this.setState(prevState => ({
                ...prevState,
                currentPanel: record,
              }));
              this.handleImportModelVisible(true);
            }}
          >
            导入
          </a>
          <Divider type="vertical" />
          <Popconfirm
            title={`确认导出米表「${record.name}」的分类与域名？`}
            onConfirm={() => {
              this.handleExport(record);
            }}
            okText="确认"
            cancelText="取消"
          >
            <a>导出</a>
          </Popconfirm>
          <Divider type="vertical" />
          <Popconfirm
            title={`确认删除米表「${record.name}」`}
            onConfirm={() => {
              this.handleDelete(record);
            }}
            okText="确认"
            cancelText="取消"
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
      type: 'panel/fetch',
      payload: formValues,
    });

    dispatch({
      type: 'panel/fetchOptions',
    });
  }

  handleDelete = (record: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'panel/remove',
      payload: record,
      callback: () => {
        dispatch({
          type: 'panel/fetch',
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
      type: 'panel/fetch',
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
      type: 'panel/fetch',
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
        type: 'panel/fetch',
        payload: values,
      });
    });
  };

  handleCreateModelVisible = (flag?: boolean) => {
    this.setState({
      createModelVisible: !!flag,
      currentRow: {},
      isEdit: false,
    });
  };

  handleImportModelVisible = (flag?: boolean) => {
    this.setState({
      importModalVisible: !!flag,
    });
  };

  handleAdd = (fields: any, isEdit: boolean) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'panel/add',
      payload: fields,
      callback: () => {
        dispatch({
          type: 'panel/fetch',
          payload: this.state.formValues,
        });
        message.success(`${isEdit ? '修改' : '添加'}成功`);
        this.handleCreateModelVisible();
      },
    });
  };

  handleImport = (fields: any) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'panel/import',
      payload: fields,
      callback: () => {
        message.success('导入成功');
        this.handleCreateModelVisible();
      },
    });
  };

  handleExport = (record: any) => {
    window.open(`/api/panel/${record.id}/export?xtoken=${localStorage.getItem('nbdomain-token')}`);
  };

  renderSimpleForm() {
    const { form } = this.props;
    const { getFieldDecorator } = form;
    return (
      <Form onSubmit={this.handleSearch} layout="inline">
        <Row gutter={{ md: 8, lg: 24, xl: 48 }}>
          <Col md={8} sm={24}>
            <FormItem label="标题">
              {getFieldDecorator('name')(<Input placeholder="请输入" />)}
            </FormItem>
          </Col>
          <Col md={8} sm={24}>
            <FormItem label="域名">
              {getFieldDecorator('domain')(<Input placeholder="请输入" />)}
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
      panel: { data, panelOptions },
      loading,
    } = this.props;

    const {
      selectedRows,
      createModelVisible,
      isEdit,
      currentPanel,
      currentRow,
      importModalVisible,
    } = this.state;

    return (
      <PageHeaderWrapper>
        <Card bordered={false}>
          <div className={styles.tableList}>
            <div className={styles.tableListForm}>{this.renderSimpleForm()}</div>
            <div className={styles.tableListOperator}>
              <Button
                icon="plus"
                type="primary"
                onClick={() => this.handleCreateModelVisible(true)}
              >
                新建
              </Button>
            </div>
            <StandardTable
              // scroll={{ x: 2400 }}
              rowKey="id"
              selectedRows={selectedRows}
              loading={loading}
              data={data}
              columns={this.columns}
              onSelectRow={this.handleSelectRows}
              onChange={this.handleStandardTableChange}
            />
          </div>
        </Card>
        <ImportForm
          panel={currentPanel}
          importModalVisible={importModalVisible}
          handleImport={this.handleImport}
          handleImportModalVisible={this.handleImportModelVisible}
        />
        <CreateForm
          handleAdd={this.handleAdd}
          handleCreateModelVisible={this.handleCreateModelVisible}
          currentRow={currentRow}
          isEdit={isEdit}
          panelOptions={panelOptions}
          createModelVisible={createModelVisible}
        />
      </PageHeaderWrapper>
    );
  }
}

export default Form.create<TableListProps>()(TableList);
