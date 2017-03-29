import React from 'react';
import { List, ListSubHeader, ListDivider } from 'react-toolbox/lib/list';
import { ListItem } from 'react-toolbox/lib/list';

class Services extends React.Component {

    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                {this.props.services.map(service => (
                    <List key={service.service} selectable ripple>
                        <ListSubHeader caption={service.service} />
                        {Object.keys(service.specs).map(spec => (
                            <ListItem
                                key={service.service + spec}
                                caption={spec}
                                legend={'Content-Type: ' + service.specs[spec].contentType}
                                onClick={() => this.props.onServiceClicked(service.service, spec, service.specs[spec])}/>
                        ))}
                        <ListDivider />
                  </List>
                ))}
            </div>
        );
    }

}

export default Services;