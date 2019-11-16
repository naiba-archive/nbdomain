import { Form, Input, Modal, Select, Upload, Button, Icon } from 'antd';

import { FormComponentProps } from 'antd/es/form';
import React from 'react';
import TextArea from 'antd/lib/input/TextArea';

const FormItem = Form.Item;

interface CreateFormProps extends FormComponentProps {
  panelOptions: any;
  modalVisible: boolean;
  handleAdd: (fieldsValue: { desc: string }) => void;
  handleModalVisible: () => void;
}
const CreateForm: React.FC<CreateFormProps> = props => {
  const { panelOptions, modalVisible, form, handleAdd, handleModalVisible } = props;
  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      form.resetFields();
      handleAdd(fieldsValue);
    });
  };
  return (
    <Modal
      destroyOnClose
      title="新建米表"
      visible={modalVisible}
      onOk={okHandle}
      onCancel={() => handleModalVisible()}
    >
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="域名">
        {form.getFieldDecorator('domain', {
          rules: [{ required: true, message: '请输入域名', min: 3 }],
        })(<Input placeholder="nai.ba" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="标题「中」">
        {form.getFieldDecorator('name', {
          rules: [{ required: true, message: '请输入标题' }],
        })(<Input placeholder="域名管理平台" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="标题「英」">
        {form.getFieldDecorator('name_en', {
          rules: [{ required: true, message: '请输入标题' }],
        })(<Input placeholder="Naiba Domain" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="简介「中」">
        {form.getFieldDecorator('desc', {
          rules: [{ required: true, message: '请输入简介' }],
        })(<TextArea placeholder="一些用爱注册的域名。" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="简介「英」">
        {form.getFieldDecorator('desc_en', {
          rules: [{ required: true, message: '请输入简介' }],
        })(<TextArea placeholder="Some domains registed by love." />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="米表主题">
        <Select defaultValue={form.getFieldValue('theme')} style={{ width: 180 }}>
          {panelOptions.themes &&
            Object.keys(panelOptions.themes).map((k: any) => (
              <Select.Option value={k}>{panelOptions.themes[k]}</Select.Option>
            ))}
        </Select>
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="报价页主题">
        <Select defaultValue={form.getFieldValue('offer_theme')} style={{ width: 180 }}>
          {panelOptions.offer_themes &&
            Object.keys(panelOptions.offer_themes).map((k: any) => (
              <Select.Option value={k}>{panelOptions.offer_themes[k]}</Select.Option>
            ))}
        </Select>
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="站点统计">
        <Select defaultValue={form.getFieldValue('analysis_type')} style={{ width: 180 }}>
          {panelOptions.analysis_types &&
            Object.keys(panelOptions.analysis_types).map((k: any) => (
              <Select.Option value={k}>{panelOptions.analysis_types[k]}</Select.Option>
            ))}
        </Select>
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="统计ID">
        {form.getFieldDecorator('analysis', {})(<Input placeholder="XA-88888" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="LOGO「中」">
        <Upload>
          <Button>
            <Icon type="upload" /> Upload
          </Button>
        </Upload>
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="LOGO「英」">
        <Upload>
          <Button>
            <Icon type="upload" /> Upload
          </Button>
        </Upload>
      </FormItem>
    </Modal>
  );
};

export default Form.create<CreateFormProps>()(CreateForm);
