import React from 'react';
import {render} from 'react-dom';

import { AppBar } from 'react-toolbox/lib/app_bar';
import { Tab, Tabs } from 'react-toolbox';

import Services from './components/services.jsx'
import Request from './components/request.jsx'

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = { services: [], tab: 0, request: {}}
        this.handleTabChange = this.handleTabChange.bind(this)
        this.fillRequest = this.fillRequest.bind(this)
    }

    handleTabChange(tab) {
        this.setState({tab});
    }

    getFakeValueFrom(fieldSpec) {
        if (typeof(fieldSpec.body) !== 'undefined') {
            return this.makeFakeBody('application/json', fieldSpec.body);
        }

        switch(fieldSpec.type) {
            case 'int': return 0;
            case 'int8': return 0;
            case 'int16': return 0;
            case 'int32': return 0;
            case 'int64': return 0;
            case 'uint': return 0;
            case 'uint8': return 0;
            case 'uint16': return 0;
            case 'uint32': return 0;
            case 'uint64': return 0;
            case 'uintptr': return 0;
            case 'byte': return 0;
            case 'rune': return 0;
            case 'float32': return 0;
            case 'float64': return 0;
            case 'bool': return false;
            default: return '';
        }
    }

    makeFakeBody(requestContentType, requestSpec) {
        if (requestContentType === 'application/json') {
            var body = {};

            for (var k in requestSpec) {
                body[k] = this.getFakeValueFrom(requestSpec[k])
            }

            return body;
        }

        return ""
    }

    fillRequest(service, procedure, spec) {
        this.setState({
            tab: 1,
            request: {
                service: service,
                procedure: procedure,
                requestContentType: spec.request.contentType,
                requestSpec: JSON.stringify(spec.request.body, null, 4),
                responseContentType: spec.response.contentType,
                responseSpec: JSON.stringify(spec.response.body, null, 4),
                fakeBody: JSON.stringify(this.makeFakeBody(spec.request.contentType, spec.request.body), null, 4)
            }
        });
    }

    render () {
        return (
            <div>
                <AppBar title='Porthos Playground'/>
                <Tabs index={this.state.tab} onChange={this.handleTabChange} fixed>
                    <Tab label='Specs'>
                        <Services services={this.state.services} onServiceClicked={this.fillRequest}/>
                    </Tab>
                    <Tab label='Request'>
                        <Request request={this.state.request}/>
                    </Tab>
                </Tabs>
            </div>
        );
    }

    componentDidMount() {
        var _this = this;

        fetch('/api/services').then(response => {
          response.json().then(function(data) {
            _this.setState({services: data})
          });
        }).catch(err => {
            console.log('Fetch Error', err);
        });
    }
}

render(<App/>, document.getElementById('app'));
