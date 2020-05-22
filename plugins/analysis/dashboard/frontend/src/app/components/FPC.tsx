import * as React from 'react';
import NodeStore from "app/stores/NodeStore";
import {inject, observer} from "mobx-react";
import {FPCStore} from "app/stores/FPCStore";
import Row from "react-bootstrap/Row";
import Container from "react-bootstrap/Container";
import Card from "react-bootstrap/Card";

interface Props {
    nodeStore?: NodeStore;
    fpcStore?: FPCStore;
}

@inject("nodeStore")
@inject("fpcStore")
@observer
export default class FPC extends React.Component<Props, any> {
    render() {
        let {conflictGrid} = this.props.fpcStore;
        return (
            <Container>
                <h3>Conflicts overview</h3>
                <Card>
                    <Card.Body>
                        <Row>
                            {conflictGrid}
                        </Row>
                    </Card.Body>
                </Card>
            </Container>
        );
    }
}
