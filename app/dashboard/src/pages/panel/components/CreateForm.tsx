import { Form, Input, Modal, Select, Upload, Button, Icon } from 'antd';

import { FormComponentProps } from 'antd/es/form';
import React from 'react';
import TextArea from 'antd/lib/input/TextArea';

const FormItem = Form.Item;

interface CreateFormProps extends FormComponentProps {
  panelOptions: any;
  isEdit: boolean;
  modalVisible: boolean;
  currentRow: any;
  handleAdd: (fieldsValue: any, isEdit: boolean) => void;
  handleModalVisible: () => void;
}

const CreateForm: React.FC<CreateFormProps> = props => {
  const {
    panelOptions,
    modalVisible,
    isEdit,
    currentRow,
    form,
    handleAdd,
    handleModalVisible,
  } = props;

  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      Object.keys(fieldsValue).forEach(k => {
        if (typeof fieldsValue[k] === 'object') {
          fieldsValue[k] = fieldsValue[k].file;
        }
      });
      handleAdd(fieldsValue, isEdit);
    });
  };

  return (
    <Modal
      destroyOnClose
      title={`${isEdit ? '修改' : '新建'}米表`}
      visible={modalVisible}
      onOk={okHandle}
      onCancel={() => handleModalVisible()}
    >
      {form.getFieldDecorator('id', {
        initialValue: currentRow.id ? currentRow.id : '',
      })(<input type="hidden" />)}
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="域名">
        {form.getFieldDecorator('domain', {
          rules: [{ required: true, message: '请输入域名', min: 3 }],
          initialValue: currentRow.domain ? currentRow.domain : '',
        })(<Input placeholder="nai.ba" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="标题「中」">
        {form.getFieldDecorator('name', {
          rules: [{ required: true, message: '请输入标题' }],
          initialValue: currentRow.name ? currentRow.name : '',
        })(<Input placeholder="域名管理平台" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="标题「英」">
        {form.getFieldDecorator('name_en', {
          rules: [{ required: true, message: '请输入标题' }],
          initialValue: currentRow.name_en ? currentRow.name_en : '',
        })(<Input placeholder="Naiba Domain" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="LOGO「中」">
        {form.getFieldDecorator('logo', {
          valuePropName: 'fileList',
          rules: [{ required: !isEdit, message: '必须上传 Logo' }],
        })(
          <Upload showUploadList={false} beforeUpload={() => false}>
            <Button>
              <Icon type="upload" /> Upload
            </Button>
          </Upload>,
        )}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="LOGO「英」">
        {form.getFieldDecorator('logo_en', {
          valuePropName: 'fileList',
          rules: [{ required: !isEdit, message: '必须上传 Logo' }],
        })(
          <Upload showUploadList={false} beforeUpload={() => false}>
            <Button>
              <Icon type="upload" /> Upload
            </Button>
          </Upload>,
        )}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="简介「中」">
        {form.getFieldDecorator('desc', {
          rules: [{ required: true, message: '请输入简介' }],
          initialValue: currentRow.desc ? currentRow.desc : '',
        })(<TextArea placeholder="一些用爱注册的域名。" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="简介「英」">
        {form.getFieldDecorator('desc_en', {
          rules: [{ required: true, message: '请输入简介' }],
          initialValue: currentRow.desc_en ? currentRow.desc_en : '',
        })(<TextArea placeholder="Some domains registed by love." />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="米表主题">
        {form.getFieldDecorator('theme', {
          rules: [{ required: true, message: '请选择一个主题' }],
          initialValue: currentRow.theme ? currentRow.theme : '',
        })(
          <Select style={{ width: 180 }}>
            {panelOptions.themes &&
              Object.keys(panelOptions.themes).map((k: any) => (
                <Select.Option key={k} value={k}>
                  {panelOptions.themes[k]}
                </Select.Option>
              ))}
          </Select>,
        )}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="报价页主题">
        {form.getFieldDecorator('offer_theme', {
          rules: [{ required: true, message: '请选择一个主题' }],
          initialValue: currentRow.offer_theme ? currentRow.offer_theme : '',
        })(
          <Select style={{ width: 180 }}>
            {panelOptions.offer_themes &&
              Object.keys(panelOptions.offer_themes).map((k: any) => (
                <Select.Option key={k} value={k}>
                  {panelOptions.offer_themes[k]}
                </Select.Option>
              ))}
          </Select>,
        )}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="站点统计">
        {form.getFieldDecorator('analysis_type', {
          initialValue: currentRow.analysis_type ? currentRow.analysis_type : '',
        })(
          <Select style={{ width: 180 }}>
            {panelOptions.analysis_types &&
              Object.keys(panelOptions.analysis_types).map((k: any) => (
                <Select.Option key={k} value={k}>
                  {panelOptions.analysis_types[k]}
                </Select.Option>
              ))}
          </Select>,
        )}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="统计ID">
        {form.getFieldDecorator('analysis', {
          initialValue: currentRow.analysis ? currentRow.analysis : '',
        })(<Input placeholder="XA-88888" />)}
      </FormItem>
    </Modal>
  );
};

export default Form.create<CreateFormProps>()(CreateForm);
