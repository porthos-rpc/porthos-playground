import React from 'react';
import Input from 'react-toolbox/lib/input';
import {Button, IconButton} from 'react-toolbox/lib/button';
import ProgressBar from 'react-toolbox/lib/progress_bar';
import Chip from 'react-toolbox/lib/chip';
import Collapsible from 'react-collapsible';
import { List, ListItem, ListDivider } from 'react-toolbox/lib/list';

import styles from "./request.css";


class Request extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            callEnabled: this.isCallEnabled(this.props.request.service, this.props.request.procedure, this.props.request.requestContentType),
            service: this.props.request.service,
            procedure: this.props.request.procedure,
            requestContentType: this.props.request.requestContentType,
            requestSpec: this.props.request.requestSpec,
            responseContentType: this.props.request.responseContentType,
            responseSpec: this.props.request.responseSpec,
            fakeBody: this.props.request.fakeBody,
            timeout: 5
        };
    }

    isCallEnabled = (service, procedure, contentType) => {
        return service && procedure && contentType
    }

    handleServiceChanged = (service) => {
        this.setState({service, callEnabled: this.isCallEnabled(service, this.state.procedure, this.state.requestContentType)})
    }

    handleProcedureChanged = (procedure) => {
        this.setState({procedure, callEnabled: this.isCallEnabled(this.state.service, procedure, this.state.requestContentType)})
    }

    handleContentTypeChanged = (contentType) => {
        this.setState({requestContentType: contentType, callEnabled: this.isCallEnabled(this.state.service, this.state.procedure, contentType)})
    }

    handleTimeoutChanged = (timeout) => {
        this.setState({timeout: parseInt(timeout)})
    }

    handleBodyChanged = (body) => {
        this.setState({body})
    }

    doRPC = () => {
        this.setState({responseError: null, response: null, loading: true})

        var _this = this
        var body = typeof(this.state.body) !== "undefined" ? this.state.body.trim() : null;

        if (body && this.state.requestContentType == 'application/json') {
            body = JSON.stringify(JSON.parse(body))
        }

        var form = new FormData();
        form.set('service', this.state.service)
        form.set('procedure', this.state.procedure)
        form.set('timeout', this.state.timeout)
        form.set('contentType', this.state.requestContentType)
        form.set('body', body)

        var payload = {
            method: "POST",
            body: form
        }

        fetch('/api/rpc', payload).then(response => {
            if (response.status === 200) {
                response.json().then(data => {
                    if (data.contentType == 'application/json') {
                        data.body = JSON.stringify(data.body, null, 4)
                    }

                    _this.setState({response: data, responseError: null, loading: false})
                }).catch(err => {
                    _this.setState({responseError: err, response: null, loading: false})
                });
            } else {
                response.text().then(data => {
                    _this.setState({responseError: data, response: null, loading: false})
                })
            }
        }).catch(err => {
            _this.setState({responseError: err, response: null, loading: false})
        });
    }

    makeObjectItem(field, fieldSpec) {
        if (typeof(fieldSpec.body) !== "undefined") {
            return this.makeSpec(field, fieldSpec.body, fieldSpec.description)
        }

        return <div>
            <div>
                <table>
                    <tbody>
                        <tr>
                            {field !== null ? <td width="300"><strong>{field}</strong></td> : null}
                            <td>{fieldSpec.type}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div className={styles.fieldDesc}>
                <span>{fieldSpec.description}</span>
            </div>
        </div>
    }

    makeObjectItems(spec) {
        return Object.keys(spec).map(field => (
            <div key={field}>
                <ListItem
                    itemContent={this.makeObjectItem(field, spec[field])}
                />
                <ListDivider/>
            </div>
        ))
    }

    makeSpec(title, spec, description) {
        if (Array.isArray(spec)) {
            return <div>
                <Collapsible trigger={title + " - array"}
                             className={styles.Collapsible}
                             openedClassName={styles.Collapsible}
                             triggerClassName={styles.Collapsible__trigger}
                             triggerOpenedClassName={styles.Collapsible__trigger}>
                    <List>
                        {this.makeObjectItem(null, spec[0])}
                    </List>
                </Collapsible>
                <div className={styles.fieldDesc}>
                    <span>{description}</span>
                </div>
            </div>
        }

        return <div>
            <Collapsible trigger={title ? title + " - object" : "object"}
                         className={styles.Collapsible}
                         openedClassName={styles.Collapsible}
                         triggerClassName={styles.Collapsible__trigger}
                         triggerOpenedClassName={styles.Collapsible__trigger}>
                <List>
                    {this.makeObjectItems(spec)}
                </List>
            </Collapsible>
            <div className={styles.fieldDesc}>
                <span>{description}</span>
            </div>
        </div>
    }

    render() {
        return (
            <div>
                <Input type='text' label='Service' name='service' defaultValue={this.props.request.service} onChange={this.handleServiceChanged}/>
                <Input type='text' label='Procedure' name='procedure' defaultValue={this.props.request.procedure} onChange={this.handleProcedureChanged}/>

                <Input type='text' label='Request Content Type' name='contentType' defaultValue={this.props.request.requestContentType} onChange={this.handleContentTypeChanged}/>
                {this.props.request.requestSpec ? this.makeSpec("Request Spec", this.props.request.requestSpec) : null}

                <Input type='text' label='Response Content Type' name='contentType' value={this.props.request.responseContentType}/>
                {this.props.request.responseSpec ? this.makeSpec("Response Spec", this.props.request.responseSpec) : null}

                <Input type='text' label='Request Body (RPC payload)' name='body' defaultValue={this.state.fakeBody} multiline onChange={this.handleBodyChanged}/>
                <Input type='number' label='Timeout' name='timeout' defaultValue={this.state.timeout.toString()} onChange={this.handleTimeoutChanged}/>
                <Button label='Call' raised primary disabled={!this.state.callEnabled} onClick={this.doRPC}/>

                <br/><br/><br/>
                {this.state.loading ? <ProgressBar mode='indeterminate'/> : null}
                {this.state.response ? <div>
                        <Input type='text' label='Status Code' name='responseStatusCode' value={this.state.response.statusCode}/>
                        <Input type='text' label='Content Type' name='responseContentType' value={this.state.response.contentType}/>
                        <Input type='text' label='Response Body' name='responseBody' value={this.state.response.body} multiline/>
                    </div>
                    : null}
                {this.state.responseError ? <Input type='text' label='ResponseError' name='responseError' value={this.state.responseError} multiline/> : null}
            </div>
        );
    }

}

export default Request;
