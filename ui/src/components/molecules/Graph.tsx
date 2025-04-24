import React from 'react';
import { Node, NodeStatus } from '../../models';
import { Step } from '../../models';
import Mermaid from '../atoms/Mermaid';

type onClickNode = (name: string) => void;

export type FlowchartType = 'TD' | 'LR';

type Props = {
  type: 'status' | 'config';
  flowchart?: FlowchartType;
  steps?: Step[] | Node[];
  onClickNode?: onClickNode;
};

declare global {
  interface Window {
    onClickMermaidNode: onClickNode;
  }
}

function Graph({
  steps,
  flowchart = 'TD',
  type = 'status',
  onClickNode,
}: Props) {
  const mermaidStyle = {
    display: 'flex',
    alignItems: 'flex-center',
    justifyContent: 'flex-start',
    width: flowchart == 'LR' && steps ? steps.length * 240 + 'px' : '100%',
    minWidth: '100%',
    minHeight: '200px',
    padding: '2em',
    borderRadius: '0.5em',
    backgroundSize: '20px 20px',
    border: '2px solid ',
  };
  var edgeId = 0;
  const graph = React.useMemo(() => {
    if (!steps) {
      return '';
    }
    const dat = flowchart == 'TD' ? ['flowchart TD;'] : ['flowchart LR;'];
    if (onClickNode) {
      window.onClickMermaidNode = onClickNode;
    }
    const addNodeFn = (step: Step, status: NodeStatus) => {
      const id = step.Name.replace(/\s/g, '_');
      const c = graphStatusMap[status] || '';
      dat.push(`${id}[${step.Name}]${c};`);

      const labelLines = step.Preconditions?.length
        ? step.Preconditions.map((c) => `${c.Condition} = ${c.Expected}`).join(
            '<br/>'
          )
        : '';

      if (step.Depends) {
        step.Depends.forEach((d) => {
          const depId = d.replace(/\s/g, '_');

          const linkOp =
            status === NodeStatus.Error ||
            status === NodeStatus.Cancel ||
            status === NodeStatus.Skipped
              ? '-.-x' // Dotted crossing for Error, Cancel, or Skipped
              : '-->'; // Default normal arrow (fallback for other statuses)
          const edgeIdStr = `e${edgeId}@`;
          const edgeStr = labelLines
            ? `${depId} ${edgeIdStr}${linkOp} |"${labelLines}"| ${id};`
            : `${depId} ${edgeIdStr}${linkOp} ${id};`;

          dat.push(edgeStr);

          if (status === NodeStatus.Running || status === NodeStatus.Success) {
            // green solid line
            dat.push(
              `linkStyle ${edgeId} stroke:lime,stroke-width:3px,color:black,background-color:black`
            );
            if (status === NodeStatus.Running) {
              dat.push(`${edgeIdStr}{ animate: true }`);
            }
          } else if (
            status === NodeStatus.Error ||
            status === NodeStatus.Cancel ||
            status === NodeStatus.Skipped
          ) {
            // red broken line
            dat.push(
              `linkStyle ${edgeId} stroke:red,stroke-width:3px,color:red`
            );
          }
        });
        edgeId++;
      }

      if (onClickNode) dat.push(`click ${id} onClickMermaidNode`);
    };

    if (type == 'status') {
      (steps as Node[]).forEach((s) => addNodeFn(s.Step, s.Status));
    } else {
      (steps as Step[]).forEach((s) => addNodeFn(s, NodeStatus.None));
    }

    dat.push(
      'linkStyle default stroke:#999,stroke-width:3px,fill:none,color:#333, '
    );
    dat.push(
      'classDef none color:#333,fill:white,stroke:lightblue,stroke-width:3px'
    );
    dat.push(
      'classDef running color:#333,fill:white,stroke:#2563eb,stroke-width:3px'
    );
    dat.push(
      'classDef error color:#333,fill:white,stroke:red,stroke-width:3px'
    );
    dat.push(
      'classDef cancel color:#333,fill:white,stroke:pink,stroke-width:3px'
    );
    dat.push(
      'classDef done color:#333,fill:white,stroke:lime,stroke-width:5px'
    );
    dat.push(
      'classDef skipped color:#333,fill:white,stroke:gray,stroke-width:3px'
    );
    return dat.join('\n');
  }, [steps, onClickNode, flowchart]);
  return <Mermaid style={mermaidStyle} def={graph} />;
}

export default Graph;

const graphStatusMap = {
  [NodeStatus.None]: ':::none',
  [NodeStatus.Running]: ':::running',
  [NodeStatus.Error]: ':::error',
  [NodeStatus.Cancel]: ':::cancel',
  [NodeStatus.Success]: ':::done',
  [NodeStatus.Skipped]: ':::skipped',
};
