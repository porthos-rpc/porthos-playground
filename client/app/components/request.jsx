import React from 'react';
import Input from 'react-toolbox/lib/input';
import {Button, IconButton} from 'react-toolbox/lib/button';
import ProgressBar from 'react-toolbox/lib/progress_bar';

class Request extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            callEnabled: this.isCallEnabled(this.props.request.service, this.props.request.procedure, this.props.request.contentType),
            service: this.props.request.service,
            procedure: this.props.request.procedure,
            contentType: this.props.request.contentType,
            timeout: 5,
            body: this.props.request.body
        };
    }

    isCallEnabled = (service, procedure, contentType) => {
        return service && procedure && contentType
    }

    handleServiceChanged = (service) => {
        this.setState({service, callEnabled: this.isCallEnabled(service, this.state.procedure, this.state.contentType)})
    }

    handleProcedureChanged = (procedure) => {
        this.setState({procedure, callEnabled: this.isCallEnabled(this.state.service, procedure, this.state.contentType)})
    }

    handleContentTypeChanged = (contentType) => {
        this.setState({contentType, callEnabled: this.isCallEnabled(this.state.service, this.state.procedure, contentType)})
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
        var body = this.state.body.trim()

        if (body && this.state.contentType == 'application/json') {
            body = JSON.stringify(JSON.parse(body))
        }

        var form = new FormData();
        form.set('service', this.state.service)
        form.set('procedure', this.state.procedure)
        form.set('contentType', this.state.contentType)
        form.set('timeout', this.state.timeout)
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

    render() {
        return (
            <div>
                <Input type='text' label='Service' name='service' defaultValue={this.props.request.service} onChange={this.handleServiceChanged}/>
                <Input type='text' label='Procedure' name='procedure' defaultValue={this.props.request.procedure} onChange={this.handleProcedureChanged}/>
                <Input type='text' label='Content Type' name='contentType' defaultValue={this.props.request.contentType} onChange={this.handleContentTypeChanged}/>
                <Input type='number' label='Timeout' name='timeout' defaultValue={this.state.timeout.toString()} onChange={this.handleTimeoutChanged}/>
                <Input type='text' label='Body' name='body' defaultValue={this.props.request.body} multiline onChange={this.handleBodyChanged}/>
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