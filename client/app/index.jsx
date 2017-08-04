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

    fillRequest(service, procedure, spec) {
        this.setState({
            tab: 1,
            request: {
                service: service,
                procedure: procedure,
                contentType: spec.contentType,
                spec: JSON.stringify(spec.body, null, 4)
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
